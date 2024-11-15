# Variáveis
ENV_FILE=".env"

# Verificar se o arquivo existe
if [[ ! -f "$ENV_FILE" ]]; then
  echo "Arquivo $ENV_FILE não encontrado."
  exit 1
fi

# Extrair variáveis do arquivo
POSTGRES_HOST=$(grep -oP '^POSTGRES_HOST=\K.*' "$ENV_FILE")
POSTGRES_PORT=$(grep -oP '^POSTGRES_PORT=\K.*' "$ENV_FILE")
POSTGRES_DB=$(grep -oP '^POSTGRES_DB=\K.*' "$ENV_FILE")
POSTGRES_USER=$(grep -oP '^POSTGRES_USER=\K.*' "$ENV_FILE")
POSTGRES_PASSWORD=$(grep -oP '^POSTGRES_PASSWORD=\K.*' "$ENV_FILE")

# Construir a URL do PostgreSQL
POSTGRES_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}"

MIGRATIONS_DIR=./database/migrations
SCHEMA_FILE=./database/schema.sql
QUERIES_DIR=./database/queries
SQLC_CONFIG=./sqlc.yaml

# Binários
MIGRATE_BIN=$(shell which migrate)
SQLC_BIN=$(shell which sqlc)

# Targets

.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  make migrate-up          - Aplica as migrações"
	@echo "  make migrate-down        - Reverte as últimas migrações"
	@echo "  make migrate-new name=XX - Cria uma nova migração (ex: name=create_users)"
	@echo "  make schema-dump         - Gera o schema.sql do banco de dados"
	@echo "  make sqlc-generate       - Gera código Go a partir das queries SQL"
	@echo "  make setup               - Executa setup completo (migração e geração de código)"

# Migrações
.PHONY: migrate-up
migrate-up:
ifndef MIGRATE_BIN
	$(error "migrate não encontrado. Instale com 'go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest'")
endif
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down 1

.PHONY: migrate-new
migrate-new:
ifndef name
	$(error "Informe o nome da migração usando 'name=nomedamigracao'")
endif
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Dump do esquema do banco
.PHONY: schema-dump
schema-dump:
	pg_dump -s -U user -d dbname > $(SCHEMA_FILE)

# Geração de código com sqlc
.PHONY: sqlc-generate
sqlc-generate:
ifndef SQLC_BIN
	$(error "sqlc não encontrado. Instale com 'go install github.com/kyleconroy/sqlc/cmd/sqlc@latest'")
endif
	$(SQLC_BIN) generate

# Setup completo
.PHONY: setup
setup: migrate-up schema-dump sqlc-generate
	@echo "Setup completo realizado com sucesso!"
