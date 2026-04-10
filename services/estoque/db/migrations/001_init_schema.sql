-- +goose Up 
CREATE SCHEMA IF NOT EXISTS estoque;

-- +goose Down 
DROP SCHEMA IF EXISTS estoque CASCADE;