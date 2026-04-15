package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/luis-octavius/Korp_Teste_Luis_Octavio/services/estoque/db/gen"
)

type EstoqueRepository interface {
	CriarProduto(ctx context.Context, codigo, nome string, saldo int32) (*db.EstoqueProduto, error)
	ListarProdutos(ctx context.Context) ([]db.EstoqueProduto, error)
	BuscarProdutoPorId(ctx context.Context, id pgtype.UUID) (*db.EstoqueProduto, error)
	DebitarEstoqueAtomico(ctx context.Context, id pgtype.UUID, quantidade int32) (*db.DebitarEstoqueAtomicoRow, error)
	RegistrarMovimentacao(ctx context.Context, params db.RegistrarMovimentacaoParams) (*db.EstoqueMovimentaco, error)
	ListarMovimentacoesProduto(ctx context.Context, produtoID pgtype.UUID) ([]db.EstoqueMovimentaco, error)
}

type estoqueRepository struct {
	queries *db.Queries
}

func NewEstoqueRepository(queries *db.Queries) EstoqueRepository {
	return &estoqueRepository{queries: queries}
}

func (r *estoqueRepository) CriarProduto(ctx context.Context, codigo, nome string, saldo int32) (*db.EstoqueProduto, error) {
	produto, err := r.queries.CriarProduto(ctx, db.CriarProdutoParams{
		Codigo: codigo,
		Nome:   nome,
		Saldo:  saldo,
	})
	if err != nil {
		return nil, fmt.Errorf("repository.CriarProduto: %w", err)
	}
	return &produto, nil
}

func (r *estoqueRepository) ListarProdutos(ctx context.Context) ([]db.EstoqueProduto, error) {
	produtos, err := r.queries.ListarProdutos(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository.ListarProdutos: %w", err)
	}

	return produtos, nil
}

func (r *estoqueRepository) BuscarProdutoPorId(ctx context.Context, id pgtype.UUID) (*db.EstoqueProduto, error) {
	produto, err := r.queries.BuscarProdutoPorId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.BuscarProdutoPorId: %w", err)
	}

	return &produto, nil
}

func (r *estoqueRepository) DebitarEstoqueAtomico(ctx context.Context, id pgtype.UUID, quantidade int32) (*db.DebitarEstoqueAtomicoRow, error) {
	produtoDebitado, err := r.queries.DebitarEstoqueAtomico(ctx, db.DebitarEstoqueAtomicoParams{
		Column1: id,
		Column2: quantidade,
	})
	if err != nil {
		return nil, fmt.Errorf("repository.DebitarEstoqueAtomico: %w", err)
	}

	return &produtoDebitado, nil
}

func (r *estoqueRepository) RegistrarMovimentacao(ctx context.Context, params db.RegistrarMovimentacaoParams) (*db.EstoqueMovimentaco, error) {
	mov, err := r.queries.RegistrarMovimentacao(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("repository.RegistrarMovimentacao: %w", err)
	}

	return &mov, nil
}

func (r *estoqueRepository) ListarMovimentacoesProduto(ctx context.Context, produtoID pgtype.UUID) ([]db.EstoqueMovimentaco, error) {
	movEstoque, err := r.queries.ListarMovimentacoesProduto(ctx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("repository.ListarMovimentacoesProduto: %w", err)
	}

	return movEstoque, nil
}
