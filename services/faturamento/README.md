# Serviço de Faturamento

Microsserviço responsável pela gestão de notas fiscais, incluindo criação,
adição de itens e impressão com débito automático de estoque.

## Tecnologias

- **Go 1.24+**
- **Chi** — roteamento HTTP
- **pgx/v5** — driver PostgreSQL
- **sqlc** — geração de queries type-safe

## Responsabilidades

- Criação de notas fiscais com numeração sequencial
- Adição de produtos e quantidades à nota
- Impressão de notas com débito automático do estoque
- Rollback automático de débitos em caso de falha

## Dependências externas

Este serviço se comunica com o **Serviço de Estoque** via HTTP para debitar
e, se necessário, reverter o saldo dos produtos no momento da impressão.

## Configuração

### Variáveis de ambiente

| Variável | Descrição | Exemplo |
|---|---|---|
| `DATABASE_URL` | String de conexão com o PostgreSQL | `postgres://korp:korp@localhost:5433/korp?sslmode=disable` |
| `PORT` | Porta do servidor HTTP | `8081` |
| `ESTOQUE_SERVICE_URL` | URL base do serviço de estoque | `http://estoque:8080` |

### Banco de dados

O serviço utiliza o schema `faturamento` dentro de um PostgreSQL compartilhado.
As tabelas são criadas automaticamente pelo `db/init.sql` na raiz do projeto.

Para rodar as migrations manualmente fora do Docker:

```bash
psql $DATABASE_URL -f db/migrations/004_faturamento_init_schema.sql
psql $DATABASE_URL -f db/migrations/005_faturamento_create_notas.sql
psql $DATABASE_URL -f db/migrations/006_faturamento_create_nota_items.sql
```

## Rodando localmente

```bash
# Instalar dependências
go mod download

# Rodar o serviço (requer o serviço de estoque rodando)
DATABASE_URL="postgres://korp:korp@localhost:5433/korp?sslmode=disable" \
PORT=8081 \
ESTOQUE_SERVICE_URL="http://localhost:8080" \
go run cmd/main.go
```

## Rodando com Docker

O modo recomendado é via Docker Compose na raiz do projeto:

```bash
docker-compose up --build
```

Para build isolado:

```bash
docker build -t korp-faturamento .
docker run -p 8081:8081 \
  -e DATABASE_URL="postgres://korp:korp@localhost:5433/korp?sslmode=disable" \
  -e ESTOQUE_SERVICE_URL="http://estoque:8080" \
  korp-faturamento
```

## Endpoints

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Health check do serviço |
| `POST` | `/notas` | Criar nova nota fiscal |
| `GET` | `/notas` | Listar notas fiscais |
| `GET` | `/notas/{id}` | Buscar nota com itens |
| `POST` | `/notas/{id}/itens` | Adicionar itens à nota |
| `POST` | `/notas/{id}/imprimir` | Imprimir nota e debitar estoque |

## Exemplos de requisição

### Health check

```bash
curl http://localhost:8081/health
```

**Resposta:**
```json
{ "status": "ok" }
```

### Criar nota fiscal

```bash
curl -X POST http://localhost:8081/notas
```

**Resposta:**
```json
{
  "id": "f3b2c4d5-...",
  "num_seq": 1,
  "status": "ABERTA"
}
```

### Adicionar itens

```bash
curl -X POST http://localhost:8081/notas/f3b2c4d5-.../itens \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"produto_id": "e2a1b3c4-...", "quantidade": 2},
      {"produto_id": "a1b2c3d4-...", "quantidade": 1}
    ]
  }'
```

**Resposta:**
```json
{ "mensagem": "itens adicionados com sucesso" }
```

### Imprimir nota

```bash
curl -X POST http://localhost:8081/notas/f3b2c4d5-.../imprimir
```

**Resposta:**
```json
{
  "id": "f3b2c4d5-...",
  "num_seq": 1,
  "status": "FECHADA",
  "printed_at": "2026-04-12T14:32:00Z"
}
```

**Erro — nota já fechada:**
```json
{ "erro": "service.ImprimirNota: nota já foi fechada ou impressa" }
```

**Erro — saldo insuficiente com rollback:**
```json
{ "erro": "service.ImprimirNota: falha ao debitar produto e2a1b3c4-...: ..." }
```

### Buscar nota com itens

```bash
curl http://localhost:8081/notas/f3b2c4d5-...
```

**Resposta:**
```json
{
  "id": "f3b2c4d5-...",
  "num_seq": 1,
  "status": "FECHADA",
  "created_at": "2026-04-12T14:30:00Z",
  "printed_at": "2026-04-12T14:32:00Z",
  "items": [
    {
      "id": "c1d2e3f4-...",
      "produto_id": "e2a1b3c4-...",
      "produto_nome": "Notebook Dell",
      "quantidade": 2
    }
  ]
}
```

## Fluxo de impressão

Ao chamar `POST /notas/{id}/imprimir`, o serviço executa:

1. Verifica se a nota está com status `ABERTA`
2. Busca todos os itens da nota
3. Para cada item, chama `POST /estoque/debitar` no serviço de estoque
4. Se qualquer débito falhar, reverte todos os débitos anteriores via `POST /estoque/reverter`
5. Fecha a nota atomicamente no banco (`status = FECHADA`)
6. Se o fechamento falhar após os débitos, executa o rollback

## Geração de código (sqlc)

Caso altere as queries em `db/queries/faturamento.sql`, regenere o código com:

```bash
make sqlc
# ou diretamente:
sqlc generate -f db/sqlc.yaml
```
