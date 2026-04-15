package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/luis-octavius/Korp_Teste_Luis_Octavio/services/faturamento/db/gen"
)

type FaturamentoRepository interface {
	CriarNota(ctx context.Context) (*db.FaturamentoNota, error)
	AdicionarItemNota(ctx context.Context, params db.AdicionarItemNotaParams) (*db.FaturamentoNotaItem, error)
	AdicionarItensNotaAtomico(ctx context.Context, itens []db.AdicionarItemNotaParams) error
	BuscarNotaPorId(ctx context.Context, id pgtype.UUID) (*db.FaturamentoNota, error)
	BuscarNotaPorNumero(ctx context.Context, numSeq int64) (*db.FaturamentoNota, error)
	ListarNotas(ctx context.Context) ([]db.FaturamentoNota, error)
	BuscarItemsNota(ctx context.Context, notaID pgtype.UUID) ([]db.FaturamentoNotaItem, error)
	FecharNotaAtomica(ctx context.Context, id pgtype.UUID) (*db.FaturamentoNota, error)
	VerificarStatusNota(ctx context.Context, id pgtype.UUID) (string, error)
}

type faturamentoRepository struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewFaturamentoRepository(queries *db.Queries, pool *pgxpool.Pool) FaturamentoRepository {
	return &faturamentoRepository{queries: queries, db: pool}
}

func (r *faturamentoRepository) CriarNota(ctx context.Context) (*db.FaturamentoNota, error) {
	nota, err := r.queries.CriarNota(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository.CriarNota: %w", err)
	}

	return &nota, nil
}

func (r *faturamentoRepository) AdicionarItemNota(ctx context.Context, args db.AdicionarItemNotaParams) (*db.FaturamentoNotaItem, error) {
	item, err := r.queries.AdicionarItemNota(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("repository.AdicionarItemNota: %w", err)
	}

	return &item, nil
}

func (r *faturamentoRepository) AdicionarItensNotaAtomico(ctx context.Context, itens []db.AdicionarItemNotaParams) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("repository.AdicionarItensNotaAtomico: erro ao iniciar transação: %w", err)
	}

	defer func() {
		if tx != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	qtx := r.queries.WithTx(tx)
	for _, item := range itens {
		if _, err := qtx.AdicionarItemNota(ctx, item); err != nil {
			return fmt.Errorf("repository.AdicionarItensNotaAtomico: erro ao inserir item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("repository.AdicionarItensNotaAtomico: erro ao commitar transação: %w", err)
	}

	tx = nil
	return nil
}

func (r *faturamentoRepository) BuscarNotaPorId(ctx context.Context, id pgtype.UUID) (*db.FaturamentoNota, error) {
	nota, err := r.queries.BuscarNotaPorId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.BuscarNotaPorId: %w", err)
	}
	return &nota, nil
}

func (r *faturamentoRepository) BuscarNotaPorNumero(ctx context.Context, numSeq int64) (*db.FaturamentoNota, error) {
	nota, err := r.queries.BuscarNotaPorNumero(ctx, numSeq)
	if err != nil {
		return nil, fmt.Errorf("repository.BuscarNotaPorNumero: %w", err)
	}
	return &nota, nil
}

func (r *faturamentoRepository) ListarNotas(ctx context.Context) ([]db.FaturamentoNota, error) {
	notas, err := r.queries.ListarNotas(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository.ListarNotas: %w", err)
	}
	return notas, nil
}

func (r *faturamentoRepository) BuscarItemsNota(ctx context.Context, notaID pgtype.UUID) ([]db.FaturamentoNotaItem, error) {
	items, err := r.queries.BuscarItemsNota(ctx, notaID)
	if err != nil {
		return nil, fmt.Errorf("repository.BuscarItemsNota: %w", err)
	}
	return items, nil
}

func (r *faturamentoRepository) FecharNotaAtomica(ctx context.Context, id pgtype.UUID) (*db.FaturamentoNota, error) {
	nota, err := r.queries.FecharNotaAtomica(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.FecharNotaAtomica: %w", err)
	}

	return &nota, nil
}

func (r *faturamentoRepository) VerificarStatusNota(ctx context.Context, id pgtype.UUID) (string, error) {
	status, err := r.queries.VerificarStatusNota(ctx, id)
	if err != nil {
		return "", fmt.Errorf("repository.VerificarStatusNota: %w", err)
	}
	return status, nil
}
