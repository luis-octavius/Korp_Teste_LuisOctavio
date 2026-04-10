-- name: CriarProduto :one 
INSERT INTO estoque.produtos (nome, saldo)
VALUES ($1, $2)
RETURNING *; 

-- name: ListarProdutos :many
SELECT * FROM estoque.produtos
ORDER BY nome;

-- name: BuscarProdutoPorId :one 
SELECT * FROM estoque.produtos 
WHERE id = $1;

-- name: DebitarEstoqueAtomico :one 
WITH debito AS (
    UPDATE estoque.produtos
    SET saldo = saldo - $2::int,
        updated_at = NOW()
    WHERE id = $1::uuid
        AND saldo >= $2::int
    RETURNING id, nome, saldo, created_at, updated_at
)
SELECT 
    id, 
    nome, 
    saldo, 
    created_at, 
    updated_at
FROM debito;

-- name: RegistrarMovimentacao :one
INSERT INTO estoque.movimentacoes (
    produto_id, 
    operacao,
    quantidade,
    motivo,
    nota_fiscal_id,
    nota_fiscal_num
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListarMovimentacoesProduto :many
SELECT * FROM estoque.movimentacoes 
WHERE produto_id = $1
ORDER BY created_at DESC;