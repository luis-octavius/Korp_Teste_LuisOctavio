package internal

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Erro string `json:"erro"`
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, mensagem string) {
	respondJSON(w, status, errorResponse{Erro: mensagem})
}
