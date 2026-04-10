-- +goose Up 
CREATE TABLE IF NOT EXISTS estoque.movimentacoes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    produto_id UUID NOT NULL REFERENCES estoque.produtos(id),
    operacao TEXT NOT NULL CHECK (operacao IN ('ENTRADA', 'SAIDA', 'AJUSTE')),
    quantidade INTEGER NOT NULL CHECK (quantidade > 0),
    motivo TEXT, 
    nota_fiscal_id UUID,
    nota_fiscal_num BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS estoque.movimentacoes;