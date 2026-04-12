-- +goose Up 
CREATE TABLE IF NOT EXISTS faturamento.nota_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nota_id UUID NOT NULL REFERENCES faturamento.notas(id) ON DELETE CASCADE,
    produto_id UUID NOT NULL,
    quantidade INTEGER NOT NULL CHECK (quantidade > 0)
);

-- +goose Down 
DROP TABLE IF EXISTS faturamento.nota_items;