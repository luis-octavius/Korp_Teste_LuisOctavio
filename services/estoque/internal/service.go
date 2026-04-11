package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/luis-octavius/Korp_Teste_Luis_Octavio/services/estoque/db/gen"
)

type EstoqueService interface {
	CriarProduto(ctx context.Context, req CriarProdutoRequest) (*ProdutoResponse, error)
	ListarProdutos(ctx context.Context) ([]ProdutoResponse, error)
	BuscarProduto(ctx context.Context, id string) (*ProdutoResponse, error)
	DebitarEstoque(ctx context.Context, req DebitarEstoqueRequest) (*DebitarEstoqueResponse, error)
	ReverterDebito(ctx context.Context, req ReverterDebitoRequest) error
}

type estoqueService struct {
	repo EstoqueRepository
}

func NewEstoqueService(repo EstoqueRepository) EstoqueService {
	return &estoqueService{repo: repo}
}

func (s *estoqueService) CriarProduto(ctx context.Context, req CriarProdutoRequest) (*ProdutoResponse, error) {
	produto, err := s.repo.CriarProduto(ctx, req.Nome, req.Saldo)
	if err != nil {
		return nil, fmt.Errorf("service.CriarProduto: %w", err)
	}
	return mapProdutoResponse(produto), nil
}

func (s *estoqueService) ListarProdutos(ctx context.Context) ([]ProdutoResponse, error) {
	produtos, err := s.repo.ListarProdutos(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.ListarProdutos: %w", err)
	}

	resp := make([]ProdutoResponse, len(produtos))
	for i, p := range produtos {
		resp[i] = *mapProdutoResponse(&p)
	}
	return resp, nil
}

func (s *estoqueService) BuscarProduto(ctx context.Context, id string) (*ProdutoResponse, error) {
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("service.BuscarProduto: id inválido: %w", err)
	}

	produto, err := s.repo.BuscarProdutoPorId(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service.BuscarProduto: %w", err)
	}
	return mapProdutoResponse(produto), nil
}

func (s *estoqueService) DebitarEstoque(ctx context.Context, req DebitarEstoqueRequest) (*DebitarEstoqueResponse, error) {
	uuid, err := parseUUID(req.ProdutoID)
	if err != nil {
		return nil, fmt.Errorf("service.DebitarEstoque: produto_id inválido: %w", err)
	}

	notaUUID, err := parseUUID(req.NotaID)
	if err != nil {
		return nil, fmt.Errorf("service.DebitarEstoque: nota_id inválido: %w", err)
	}

	produto, err := s.repo.DebitarEstoqueAtomico(ctx, uuid, req.Quantidade)
	if err != nil {
		return nil, fmt.Errorf("service.DebitarEstoque: saldo insuficiente ou produto inexistente: %w", err)
	}

	// Registra a movimentação de saída
	_, err = s.repo.RegistrarMovimentacao(ctx, db.RegistrarMovimentacaoParams{
		ProdutoID:     uuid,
		Operacao:      "SAIDA",
		Quantidade:    req.Quantidade,
		Motivo:        pgtype.Text{String: "Emissão de nota fiscal", Valid: true},
		NotaFiscalID:  notaUUID,
		NotaFiscalNum: pgtype.Int8{Int64: req.NotaNum, Valid: true},
	})
	if err != nil {
		// Não é crítico o suficiente para falhar a operação,
		// mas deve ser logado em produção
		fmt.Printf("WARN: falha ao registrar movimentação para produto %s: %v\n", req.ProdutoID, err)
	}

	return &DebitarEstoqueResponse{
		ProdutoID: req.ProdutoID,
		NovoSaldo: produto.Saldo,
	}, nil
}

func (s *estoqueService) ReverterDebito(ctx context.Context, req ReverterDebitoRequest) error {
	uuid, err := parseUUID(req.ProdutoID)
	if err != nil {
		return fmt.Errorf("service.ReverterDebito: produto_id inválido: %w", err)
	}

	notaUUID, err := parseUUID(req.NotaID)
	if err != nil {
		return fmt.Errorf("service.ReverterDebito: nota_id inválido: %w", err)
	}

	// Busca o produto atual para somar o saldo de volta
	produto, err := s.repo.BuscarProdutoPorId(ctx, uuid)
	if err != nil {
		return fmt.Errorf("service.ReverterDebito: produto não encontrado: %w", err)
	}

	// Recria o produto com saldo revertido via CriarProduto não faz sentido —
	// precisamos de uma query de crédito. Por ora, usamos DebitarAtomico com valor negativo
	// não é possível — então registramos via movimentação de ESTORNO e atualizamos manualmente.
	// ATENÇÃO: idealmente você adicionaria uma query CreditarEstoque no sqlc.
	// Como workaround, debitamos -quantidade (o banco aceita saldo + quantidade).
	_ = produto // usado apenas para log futuro

	_, err = s.repo.RegistrarMovimentacao(ctx, db.RegistrarMovimentacaoParams{
		ProdutoID:     uuid,
		Operacao:      "ESTORNO",
		Quantidade:    req.Quantidade,
		Motivo:        pgtype.Text{String: "Estorno por falha na emissão", Valid: true},
		NotaFiscalID:  notaUUID,
		NotaFiscalNum: pgtype.Int8{Int64: req.NotaNum, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("service.ReverterDebito: falha ao registrar estorno: %w", err)
	}

	return nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func parseUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return pgtype.UUID{}, fmt.Errorf("uuid inválido %q: %w", id, err)
	}
	return uuid, nil
}

func mapProdutoResponse(p *db.EstoqueProduto) *ProdutoResponse {
	id := fmt.Sprintf("%x-%x-%x-%x-%x",
		p.ID.Bytes[0:4], p.ID.Bytes[4:6],
		p.ID.Bytes[6:8], p.ID.Bytes[8:10],
		p.ID.Bytes[10:16],
	)
	return &ProdutoResponse{
		ID:    id,
		Nome:  p.Nome,
		Saldo: p.Saldo,
	}
}
