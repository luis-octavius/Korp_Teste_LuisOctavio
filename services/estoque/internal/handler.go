package internal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type EstoqueHandler struct {
	service EstoqueService
}

func NewEstoqueHandler(service EstoqueService) *EstoqueHandler {
	return &EstoqueHandler{service: service}
}

func (h *EstoqueHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)

	r.Route("/produtos", func(r chi.Router) {
		r.Post("/", h.CriarProduto)
		r.Get("/", h.ListarProdutos)
		r.Get("/{id}", h.BuscarProduto)
	})

	r.Route("/estoque", func(r chi.Router) {
		r.Post("/debitar", h.DebitarEstoque)
		r.Post("/reverter", h.ReverterDebito)
	})
}

// GET /health
func (h *EstoqueHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /produtos
func (h *EstoqueHandler) CriarProduto(w http.ResponseWriter, r *http.Request) {
	var req CriarProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.Nome == "" {
		respondError(w, http.StatusBadRequest, "campo 'nome' é obrigatório")
		return
	}
	if req.Saldo < 0 {
		respondError(w, http.StatusBadRequest, "campo 'saldo' não pode ser negativo")
		return
	}

	produto, err := h.service.CriarProduto(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, produto)
}

// GET /produtos
func (h *EstoqueHandler) ListarProdutos(w http.ResponseWriter, r *http.Request) {
	produtos, err := h.service.ListarProdutos(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, produtos)
}

// GET /produtos/{id}
func (h *EstoqueHandler) BuscarProduto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	produto, err := h.service.BuscarProduto(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "produto não encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, produto)
}

// POST /estoque/debitar
func (h *EstoqueHandler) DebitarEstoque(w http.ResponseWriter, r *http.Request) {
	var req DebitarEstoqueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.ProdutoID == "" {
		respondError(w, http.StatusBadRequest, "campo 'produto_id' é obrigatório")
		return
	}
	if req.Quantidade <= 0 {
		respondError(w, http.StatusBadRequest, "campo 'quantidade' deve ser maior que zero")
		return
	}

	resp, err := h.service.DebitarEstoque(r.Context(), req)
	if err != nil {
		// Saldo insuficiente é um erro de negócio, não erro interno
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// POST /estoque/reverter
func (h *EstoqueHandler) ReverterDebito(w http.ResponseWriter, r *http.Request) {
	var req ReverterDebitoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.ProdutoID == "" {
		respondError(w, http.StatusBadRequest, "campo 'produto_id' é obrigatório")
		return
	}
	if req.Quantidade <= 0 {
		respondError(w, http.StatusBadRequest, "campo 'quantidade' deve ser maior que zero")
		return
	}

	if err := h.service.ReverterDebito(r.Context(), req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"mensagem": "estorno realizado com sucesso"})
}
