package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type FaturamentoHandler struct {
	service FaturamentoService
}

func NewFaturamentoHandler(service FaturamentoService) *FaturamentoHandler {
	return &FaturamentoHandler{service: service}
}

// RegisterRoutes define as rotas do serviço de faturamento,
// incluindo endpoints para criar notas, listar notas, buscar nota por ID,
// adicionar itens a uma nota e imprimir a nota.
func (h *FaturamentoHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)

	r.Route("/notas", func(r chi.Router) {
		r.Post("/", h.CriarNota)
		r.Get("/", h.ListarNotas)
		r.Get("/{id}", h.BuscarNota)
		r.Post("/{id}/itens", h.AdicionarItens)
		r.Post("/{id}/imprimir", h.ImprimirNota)
	})
}

// GET /health
func (h *FaturamentoHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /notas
func (h *FaturamentoHandler) CriarNota(w http.ResponseWriter, r *http.Request) {
	nota, err := h.service.CriarNota(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, nota)
}

// GET /notas
func (h *FaturamentoHandler) ListarNotas(w http.ResponseWriter, r *http.Request) {
	notas, err := h.service.ListarNotas(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, notas)
}

// GET /notas/{id}
// Endpoint para buscar uma nota por ID,
// utilizado durante o processo de faturamento para
// consultar o status da nota e seus detalhes.
func (h *FaturamentoHandler) BuscarNota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	nota, err := h.service.BuscarNota(r.Context(), id)
	if err != nil {
		if isInputValidationError(err) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "nota não encontrada")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nota)
}

// POST /notas/{id}/itens
// Endpoint para adicionar itens a uma nota existente,
// utilizado durante o processo de faturamento para compor a nota
// com os produtos e quantidades desejados.
func (h *FaturamentoHandler) AdicionarItens(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	var req AdicionarItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if len(req.Items) == 0 {
		respondError(w, http.StatusBadRequest, "a nota deve conter ao menos um item")
		return
	}

	req.NotaID = id

	if err := h.service.AdicionarItens(r.Context(), req); err != nil {
		if isInputValidationError(err) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "nota não encontrada")
			return
		}
		// Nota fechada é erro de negócio
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"mensagem": "itens adicionados com sucesso"})
}

// POST /notas/{id}/imprimir
// Endpoint para fechar a nota e gerar o documento fiscal,
func (h *FaturamentoHandler) ImprimirNota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	nota, err := h.service.ImprimirNota(r.Context(), id)
	if err != nil {
		if isInputValidationError(err) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "nota não encontrada")
			return
		}
		// Nota já fechada ou saldo insuficiente são erros de negócio
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nota)
}

// isInputValidationError verifica se o erro é relacionado
// a validação de entrada, como campos obrigatórios ou formatos inválidos.
func isInputValidationError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "inválido") ||
		strings.Contains(msg, "invalido") ||
		strings.Contains(msg, "deve ser maior que zero")
}
