-- +goose Up 
CREATE SEQUENCE IF NOT EXISTS faturamento.nota_num_seq;

CREATE TABLE IF NOT EXISTS faturamento.notas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    num_seq BIGINT NOT NULL DEFAULT nextval('faturamento.nota_num_seq'),
    status TEXT NOT NULL DEFAULT 'ABERTA' CHECK (status IN ('ABERTA', 'FECHADA')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    printed_at TIMESTAMPTZ,
    UNIQUE (num_seq)
);

-- +goose Down
DROP TABLE IF EXISTS faturamento.notas;
DROP SEQUENCE IF EXISTS faturamento.nota_num_seq;