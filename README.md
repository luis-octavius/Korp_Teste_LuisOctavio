# Korp_Teste_Luis_Octavio

Sistema de emissão de Notas Fiscais desenvolvido como teste técnico para a Korp ERP (Viasoft).

## Arquitetura  

```text
O sistema é composto por dois microsserviços independentes que se comunicam via HTTP,
compartilhando um único banco de dados PostgreSQL com schemas isolados.  
frontend/          → Aplicação Angular  
services/  
estoque/         → Controle de produtos e saldos (porta 8080)  
faturamento/     → Gestão de notas fiscais (porta 8081)  
db/  
init.sql         → Script de inicialização do banco  
docker-compose.yml → Orquestração dos serviços 
```

## Tecnologias  

| Camada | Tecnologia |
|---|---|
| Frontend | Angular |
| Backend | Go 1.24+ |
| Roteamento HTTP | Chi |
| Driver PostgreSQL | pgx/v5 |
| Geração de queries | sqlc |
| Banco de dados | PostgreSQL 16 |
| Containerização | Docker + Docker Compose |

## Pré-requisitos  

- Docker
- Docker Compose

## Subindo o projeto completo  

```bash
# Primeira execução
docker-compose up --build

# Execuções seguintes
docker-compose up

# Derrubar os serviços
docker-compose down

# Derrubar e limpar o banco de dados
docker-compose down -v
```

## Serviços  

| Serviço | URL local | Documentação |
|---|---|---|
| Estoque | http://localhost:8080 | [services/estoque/README.md](services/estoque/README.md) |
| Faturamento | http://localhost:8081 | [services/faturamento/README.md](services/faturamento/README.md) |

## Banco de dados  

Um único PostgreSQL com dois schemas isolados:

| Schema | Responsável |
|---|---|
| `estoque` | Serviço de Estoque |
| `faturamento` | Serviço de Faturamento |

O script `db/init.sql` cria todos os schemas e tabelas automaticamente
na primeira inicialização do container.

## Fluxo principal
```text
1 - Cadastrar produtos         → POST /produtos          (estoque:8080)  
2 - Criar nota fiscal          → POST /notas             (faturamento:8081)  
3 - Adicionar itens à nota     → POST /notas/{id}/itens  (faturamento:8081)  
4 - Imprimir nota              → POST /notas/{id}/imprimir  
↓  
Debita estoque via HTTP para cada item  
↓  
Se falhar → rollback automático dos débitos  
↓  
Fecha a nota atomicamente (status = FECHADA)  
```

## Requisitos implementados

- [x] Arquitetura de microsserviços
- [x] Cadastro de produtos com controle de saldo
- [x] Criação de notas fiscais com numeração sequencial
- [x] Impressão de notas com débito automático de estoque
- [x] Rollback automático em caso de falha
- [x] Persistência em banco de dados relacional
- [x] Health checks em todos os serviços
- [x] Graceful shutdown
