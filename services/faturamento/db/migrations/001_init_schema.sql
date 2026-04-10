-- +goose Up 
CREATE SCHEMA IF NOT EXISTS faturamento; 

-- +goose Down 
DROP SCHEMA IF EXISTS faturamento;