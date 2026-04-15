package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/luis-octavius/Korp_Teste_Luis_Octavio/services/faturamento/db/gen"
	internaldb "github.com/luis-octavius/Korp_Teste_Luis_Octavio/services/faturamento/internal/db"
)

type App struct {
	router *chi.Mux
	db     *pgxpool.Pool
}

func NewApp(ctx context.Context) (*App, error) {
	// 1. Conexão com o banco
	pool, err := internaldb.NovaConexao(ctx)
	if err != nil {
		return nil, fmt.Errorf("app: falha ao conectar ao banco: %w", err)
	}

	// 2. Wire de dependências
	queries := db.New(pool)
	repo := NewFaturamentoRepository(queries, pool)
	service := NewFaturamentoService(repo)
	handler := NewFaturamentoHandler(service)

	// 3. Router
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	handler.RegisterRoutes(r)

	return &App{router: r, db: pool}, nil
}

func (a *App) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("faturamento: servidor rodando na porta %s\n", port)
	return http.ListenAndServe(":"+port, a.router)
}

func (a *App) Shutdown() {
	a.db.Close()
}
