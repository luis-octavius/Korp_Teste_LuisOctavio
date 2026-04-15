-- +goose Up
ALTER TABLE estoque.produtos
ADD COLUMN IF NOT EXISTS codigo TEXT;

UPDATE estoque.produtos
SET codigo = id::text
WHERE codigo IS NULL OR btrim(codigo) = '';

ALTER TABLE estoque.produtos
ALTER COLUMN codigo SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_estoque_produtos_codigo
ON estoque.produtos (codigo);

-- +goose Down
DROP INDEX IF EXISTS idx_estoque_produtos_codigo;

ALTER TABLE estoque.produtos
DROP COLUMN IF EXISTS codigo;
