-- ─── Estoque ─────────────────────────────────────────────────────────────────

CREATE SCHEMA IF NOT EXISTS estoque;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS estoque.produtos (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    nome       TEXT        NOT NULL,
    saldo      INTEGER     NOT NULL DEFAULT 0 CHECK (saldo >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS estoque.movimentacoes (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    produto_id      UUID        NOT NULL REFERENCES estoque.produtos(id),
    operacao        TEXT        NOT NULL CHECK (operacao IN ('ENTRADA', 'SAIDA', 'AJUSTE', 'ESTORNO')),
    quantidade      INTEGER     NOT NULL CHECK (quantidade > 0),
    motivo          TEXT,
    nota_fiscal_id  UUID,
    nota_fiscal_num BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ─── Faturamento ─────────────────────────────────────────────────────────────

CREATE SCHEMA IF NOT EXISTS faturamento;

CREATE SEQUENCE IF NOT EXISTS faturamento.nota_num_seq;

CREATE TABLE IF NOT EXISTS faturamento.notas (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    num_seq    BIGINT      NOT NULL DEFAULT nextval('faturamento.nota_num_seq'),
    status     TEXT        NOT NULL DEFAULT 'ABERTA' CHECK (status IN ('ABERTA', 'FECHADA')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    printed_at TIMESTAMPTZ,
    UNIQUE (num_seq)
);

CREATE TABLE IF NOT EXISTS faturamento.nota_items (
    id         UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    nota_id    UUID    NOT NULL REFERENCES faturamento.notas(id) ON DELETE CASCADE,
    produto_id UUID    NOT NULL,
    quantidade INTEGER NOT NULL CHECK (quantidade > 0)
);
