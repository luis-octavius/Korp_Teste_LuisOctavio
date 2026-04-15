# Serviço de Estoque

Microsserviço responsável pelo controle de produtos e saldos em estoque.

## Tecnologias

- **Go 1.24+**
- **Chi** — roteamento HTTP
- **pgx/v5** — driver PostgreSQL
- **sqlc** — geração de queries type-safe

## Responsabilidades

- Cadastro e consulta de produtos
- Débito atômico de estoque via `UPDATE ... WHERE saldo >= quantidade`
- Reversão de débito (estorno) em caso de falha no faturamento
- Registro de movimentações (`ENTRADA`, `SAIDA`, `AJUSTE`, `ESTORNO`)

## Configuração

### Variáveis de ambiente

| Variável | Descrição | Exemplo |
|---|---|---|
| `DATABASE_URL` | String de conexão com o PostgreSQL | `postgres://korp:korp@localhost:5433/korp?sslmode=disable` |
| `PORT` | Porta do servidor HTTP | `8080` |

### Banco de dados

O serviço utiliza o schema `estoque` dentro de um PostgreSQL compartilhado.
As tabelas são criadas automaticamente pelo `db/init.sql` na raiz do projeto.

Para rodar as migrations manualmente fora do Docker:

```bash
psql $DATABASE_URL -f db/migrations/001_estoque_init_schema.sql
psql $DATABASE_URL -f db/migrations/002_estoque_create_produtos.sql
psql $DATABASE_URL -f db/migrations/003_estoque_create_movimentacoes.sql
```

## Rodando localmente

```bash
# Instalar dependências
go mod download

# Rodar o serviço
DATABASE_URL="postgres://korp:korp@localhost:5433/korp?sslmode=disable" \
PORT=8080 \
go run cmd/main.go
```

## Rodando com Docker

O modo recomendado é via Docker Compose na raiz do projeto:

```bash
docker-compose up --build
```

Para build isolado:

```bash
docker build -t korp-estoque .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://korp:korp@localhost:5433/korp?sslmode=disable" \
  korp-estoque
```

## Endpoints

### Produtos

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Health check do serviço |
| `POST` | `/produtos` | Cadastrar produto |
| `GET` | `/produtos` | Listar produtos |
| `GET` | `/produtos/{id}` | Buscar produto por ID |

### Estoque (uso interno)

> Estes endpoints são consumidos pelo serviço de Faturamento e não devem ser expostos publicamente.

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/estoque/debitar` | Debitar quantidade do estoque |
| `POST` | `/estoque/reverter` | Reverter débito (estorno) |

## Exemplos de requisição

### Health check

```bash
curl http://localhost:8080/health
```

**Resposta:**
```json
{ "status": "ok" }
```

### Cadastrar produto

```bash
curl -X POST http://localhost:8080/produtos \
  -H "Content-Type: application/json" \
  -d '{"codigo": "SKU-001", "nome": "Notebook Dell", "saldo": 10}'
```

**Resposta:**
```json
{
  "id": "e2a1b3c4-...",
  "nome": "Notebook Dell",
  "saldo": 10
}
```

### Listar produtos

```bash
curl http://localhost:8080/produtos
```

### Buscar produto por ID

```bash
curl http://localhost:8080/produtos/e2a1b3c4-...
```

### Debitar estoque

```bash
curl -X POST http://localhost:8080/estoque/debitar \
  -H "Content-Type: application/json" \
  -d '{
    "produto_id": "e2a1b3c4-...",
    "quantidade": 2,
    "nota_id": "f3b2c4d5-...",
    "nota_num": 1
  }'
```

**Resposta:**
```json
{
  "produto_id": "e2a1b3c4-...",
  "novo_saldo": 8
}
```

**Erro — saldo insuficiente:**
```json
{
  "erro": "service.DebitarEstoque: saldo insuficiente ou produto inexistente"
}
```

### Reverter débito

```bash
curl -X POST http://localhost:8080/estoque/reverter \
  -H "Content-Type: application/json" \
  -d '{
    "produto_id": "e2a1b3c4-...",
    "quantidade": 2,
    "nota_id": "f3b2c4d5-...",
    "nota_num": 1
  }'
```

**Resposta:**
```json
{ "mensagem": "estorno realizado com sucesso" }
```

## Geração de código (sqlc)

Caso altere as queries em `db/queries/estoque.sql`, regenere o código com:

```bash
make sqlc
# ou diretamente:
sqlc generate -f db/sqlc.yaml
```
