package internal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type FaturamentoHandler struct {
	service FaturamentoService
}

func NewFaturamentoHandler(service FaturamentoService) *FaturamentoHandler {
	return &FaturamentoHandler{service: service}
}

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
func (h *FaturamentoHandler) BuscarNota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	nota, err := h.service.BuscarNota(r.Context(), id)
	if err != nil {
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
		// Nota fechada é erro de negócio
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"mensagem": "itens adicionados com sucesso"})
}

// POST /notas/{id}/imprimir
func (h *FaturamentoHandler) ImprimirNota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	nota, err := h.service.ImprimirNota(r.Context(), id)
	if err != nil {
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
