package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/luis-octavius/Korp_Teste_Luis_Octavio/services/faturamento/db/gen"
)

type FaturamentoService interface {
	CriarNota(ctx context.Context) (*CriarNotaResponse, error)
	AdicionarItens(ctx context.Context, req AdicionarItemRequest) error
	BuscarNota(ctx context.Context, id string) (*NotaDetalheResponse, error)
	ListarNotas(ctx context.Context) ([]CriarNotaResponse, error)
	ImprimirNota(ctx context.Context, id string) (*ImprimirNotaResponse, error)
}

type faturamentoService struct {
	repo       FaturamentoRepository
	estoqueURL string
	httpClient *http.Client
}

func NewFaturamentoService(repo FaturamentoRepository) FaturamentoService {
	estoqueURL := os.Getenv("ESTOQUE_SERVICE_URL")
	if estoqueURL == "" {
		estoqueURL = "http://estoque:8080"
	}

	return &faturamentoService{
		repo:       repo,
		estoqueURL: estoqueURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *faturamentoService) CriarNota(ctx context.Context) (*CriarNotaResponse, error) {
	nota, err := s.repo.CriarNota(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.CriarNota: %w", err)
	}
	return &CriarNotaResponse{
		ID:     uuidToString(nota.ID.Bytes),
		NumSeq: nota.NumSeq,
		Status: nota.Status,
	}, nil
}

func (s *faturamentoService) AdicionarItens(ctx context.Context, req AdicionarItemRequest) error {
	notaUUID, err := parseUUID(req.NotaID)
	if err != nil {
		return fmt.Errorf("service.AdicionarItens: nota_id inválido: %w", err)
	}

	// Verifica se a nota está aberta antes de adicionar itens
	status, err := s.repo.VerificarStatusNota(ctx, notaUUID)
	if err != nil {
		return fmt.Errorf("service.AdicionarItens: nota não encontrada: %w", err)
	}
	if status != "ABERTA" {
		return fmt.Errorf("service.AdicionarItens: nota %s não está aberta", req.NotaID)
	}

	for _, item := range req.Items {
		produtoUUID, err := parseUUID(item.ProdutoID)
		if err != nil {
			return fmt.Errorf("service.AdicionarItens: produto_id inválido %s: %w", item.ProdutoID, err)
		}

		_, err = s.repo.AdicionarItemNota(ctx, db.AdicionarItemNotaParams{
			NotaID:     notaUUID,
			ProdutoID:  produtoUUID,
			Quantidade: item.Quantidade,
		})
		if err != nil {
			return fmt.Errorf("service.AdicionarItens: erro ao adicionar item %s: %w", item.ProdutoID, err)
		}
	}
	return nil
}

func (s *faturamentoService) BuscarNota(ctx context.Context, id string) (*NotaDetalheResponse, error) {
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("service.BuscarNota: id inválido: %w", err)
	}

	nota, err := s.repo.BuscarNotaPorId(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service.BuscarNota: %w", err)
	}

	items, err := s.repo.BuscarItemsNota(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service.BuscarNota: erro ao buscar itens: %w", err)
	}

	resp := &NotaDetalheResponse{
		ID:        uuidToString(nota.ID.Bytes),
		NumSeq:    nota.NumSeq,
		Status:    nota.Status,
		CreatedAt: nota.CreatedAt.Time.Format(time.RFC3339),
		Items:     make([]ItemNotaResponse, len(items)),
	}

	if nota.PrintedAt.Valid {
		t := nota.PrintedAt.Time.Format(time.RFC3339)
		resp.PrintedAt = &t
	}

	cache := make(map[string]*EstoqueProdutoResponse)

	for i, item := range items {
		produtoID := uuidToString(item.ProdutoID.Bytes)

		produto, ok := cache[produtoID]
		if !ok {
			p, err := s.buscarProduto(ctx, produtoID)
			if err != nil {
				fmt.Printf("WARN: service.BuscarNota: não foi possível buscar nome do produto %s: %v\n", produtoID, err)
			} else {
				cache[produtoID] = p
				produto = p
			}
		}

		nomeProduto := ""
		if produto != nil {
			nomeProduto = produto.Nome
		}

		resp.Items[i] = ItemNotaResponse{
			ID:          uuidToString(item.ID.Bytes),
			ProdutoID:   uuidToString(item.ProdutoID.Bytes),
			Quantidade:  item.Quantidade,
			ProdutoNome: nomeProduto,
		}
	}

	return resp, nil
}

func (s *faturamentoService) ListarNotas(ctx context.Context) ([]CriarNotaResponse, error) {
	notas, err := s.repo.ListarNotas(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.ListarNotas: %w", err)
	}

	resp := make([]CriarNotaResponse, len(notas))
	for i, n := range notas {
		resp[i] = CriarNotaResponse{
			ID:     uuidToString(n.ID.Bytes),
			NumSeq: n.NumSeq,
			Status: n.Status,
		}
	}
	return resp, nil
}

// ImprimirNota é o método mais crítico: debita o estoque item a item
// e faz rollback manual se qualquer débito falhar.
func (s *faturamentoService) ImprimirNota(ctx context.Context, id string) (*ImprimirNotaResponse, error) {
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("service.ImprimirNota: id inválido: %w", err)
	}

	// 1. Verifica se a nota está aberta
	status, err := s.repo.VerificarStatusNota(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service.ImprimirNota: nota não encontrada: %w", err)
	}
	if status != "ABERTA" {
		return nil, fmt.Errorf("service.ImprimirNota: nota já foi fechada ou impressa")
	}

	// 2. Busca a nota completa para ter o num_seq
	nota, err := s.repo.BuscarNotaPorId(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service.ImprimirNota: %w", err)
	}

	// 3. Busca os itens da nota
	items, err := s.repo.BuscarItemsNota(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service.ImprimirNota: erro ao buscar itens: %w", err)
	}

	// 4. Debita o estoque item a item, guardando os já debitados para rollback
	debitados := make([]DebitarEstoqueRequest, 0, len(items))

	for _, item := range items {
		req := DebitarEstoqueRequest{
			ProdutoID:  uuidToString(item.ProdutoID.Bytes),
			Quantidade: item.Quantidade,
			NotaID:     id,
			NotaNum:    nota.NumSeq,
		}

		if err := s.debitarEstoque(ctx, req); err != nil {
			// Falhou — inicia rollback de todos os itens já debitados
			s.rollbackDebitos(ctx, debitados, id, nota.NumSeq)
			return nil, fmt.Errorf("service.ImprimirNota: falha ao debitar produto %s: %w", req.ProdutoID, err)
		}

		debitados = append(debitados, req)
	}

	// 5. Fecha a nota atomicamente no banco
	notaFechada, err := s.repo.FecharNotaAtomica(ctx, uuid)
	if err != nil {
		// Débitos já foram feitos mas a nota não fechou — faz rollback
		s.rollbackDebitos(ctx, debitados, id, nota.NumSeq)
		return nil, fmt.Errorf("service.ImprimirNota: falha ao fechar nota: %w", err)
	}

	return &ImprimirNotaResponse{
		ID:        uuidToString(notaFechada.ID.Bytes),
		NumSeq:    notaFechada.NumSeq,
		Status:    notaFechada.Status,
		PrintedAt: notaFechada.PrintedAt.Time.Format(time.RFC3339),
	}, nil
}

// ─── Comunicação com o serviço de estoque ──────────────────────────────────

func (s *faturamentoService) debitarEstoque(ctx context.Context, req DebitarEstoqueRequest) error {
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.estoqueURL+"/estoque/debitar",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("debitarEstoque: erro ao criar request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("debitarEstoque: serviço de estoque indisponível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("debitarEstoque: estoque retornou status %d para produto %s", resp.StatusCode, req.ProdutoID)
	}
	return nil
}

func (s *faturamentoService) rollbackDebitos(
	ctx context.Context,
	debitados []DebitarEstoqueRequest,
	notaID string,
	notaNum int64,
) {
	for _, d := range debitados {
		req := ReverterDebitoRequest{
			ProdutoID:  d.ProdutoID,
			Quantidade: d.Quantidade,
			NotaID:     notaID,
			NotaNum:    notaNum,
		}

		body, _ := json.Marshal(req)
		httpReq, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			s.estoqueURL+"/estoque/reverter",
			bytes.NewReader(body),
		)
		if err != nil {
			fmt.Printf("WARN: rollback falhou ao criar request para produto %s: %v\n", d.ProdutoID, err)
			continue
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(httpReq)
		if err != nil {
			fmt.Printf("WARN: rollback falhou para produto %s: %v\n", d.ProdutoID, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("WARN: rollback retornou status %d para produto %s\n", resp.StatusCode, d.ProdutoID)
		}
	}
}

// ─── Estoque ───────────────────────────────────────────────────────────────
func (s *faturamentoService) buscarProduto(ctx context.Context, produtoID string) (*EstoqueProdutoResponse, error) {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.estoqueURL+"/produtos/"+produtoID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("buscarProduto: erro ao criar request: %w", err)
	}

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("buscarProduto: serviço de estoque indisponível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("buscarProduto: estoque retornou status %s", resp.Status)
	}

	if resp.StatusCode == http.StatusOK {
		return nil, fmt.Errorf("buscarProduto: estoque retornou status %s", resp.Status)
	}

	var produto EstoqueProdutoResponse
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return nil, fmt.Errorf("buscarProduto: erro ao decodificar resposta: %w", err)
	}

	return &produto, nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func uuidToString(b [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	)
}

func parseUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return pgtype.UUID{}, fmt.Errorf("uuid inválido %q: %w", id, err)
	}
	return uuid, nil
}
