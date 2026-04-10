DB_URL ?= postgres://postgres:postgres@localhost:5432/korp?sslmode=disable
MIGRATIONS_DIR := db/migrations 

.PHONY: migrate-up migrate-down migrate-status migrate-new sqlc

migrate-up: goose -dir ($(MIGRATIONS_DIR)) postgres "$(DB_URL)" up
migrate-down: goose -dir ($(MIGRATIONS_DIR)) postgres "$(DB_URL)" down
migrate-status: goose -dir ($(MIGRATIONS_DIR)) postgres "$(DB_URL)" status
migrate-new: @echo "Use: make migrate-new 
NAME=create_xxx" goose-dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" create $(NAME) sql

sqlc:
	sqlc generate -f sqlc.yaml