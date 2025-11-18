.PHONY: test coverage test-verbose test-utils test-service help

# Variáveis
GO_VERSION := 1.23-alpine
DOCKER_IMAGE := golang:$(GO_VERSION)
DOCKER_RUN := docker run --rm -v "$(PWD)":/app -w /app $(DOCKER_IMAGE)

# Comando padrão
.DEFAULT_GOAL := help

## help: Mostra esta mensagem de ajuda
help:
	@echo "Comandos disponíveis:"
	@echo ""
	@echo "  make test          - Executa todos os testes"
	@echo "  make coverage      - Executa testes com cobertura"
	@echo "  make test-verbose  - Executa testes com output detalhado"
	@echo "  make test-utils    - Executa apenas testes do pacote utils"
	@echo "  make test-service  - Executa apenas testes do pacote service"
	@echo "  make build         - Compila a aplicação"
	@echo "  make run           - Executa a aplicação localmente (requer Go instalado)"
	@echo "  make docker-up     - Inicia os containers Docker"
	@echo "  make docker-down   - Para os containers Docker"
	@echo "  make docker-logs   - Mostra logs dos containers"
	@echo ""

## test: Executa todos os testes unitários
test:
	@echo "🧪 Executando testes..."
	@$(DOCKER_RUN) go test ./...

## coverage: Executa testes com cobertura de código
coverage:
	@echo "📊 Executando testes com cobertura..."
	@$(DOCKER_RUN) go test -cover ./...

## test-verbose: Executa testes com output detalhado
test-verbose:
	@echo "🔍 Executando testes com output detalhado..."
	@$(DOCKER_RUN) go test -v ./...

## test-utils: Executa apenas testes do pacote utils
test-utils:
	@echo "🧪 Executando testes do pacote utils..."
	@$(DOCKER_RUN) go test -v ./internal/utils

## test-service: Executa apenas testes do pacote service
test-service:
	@echo "🧪 Executando testes do pacote service..."
	@$(DOCKER_RUN) go test -v ./internal/service

## build: Compila a aplicação
build:
	@echo "🔨 Compilando aplicação..."
	@$(DOCKER_RUN) go build -o bin/api ./cmd/api

## run: Executa a aplicação localmente (requer Go instalado)
run:
	@echo "🚀 Executando aplicação..."
	@go run cmd/api/main.go

## docker-up: Inicia os containers Docker
docker-up:
	@echo "🐳 Iniciando containers Docker..."
	@docker compose up -d

## docker-down: Para os containers Docker
docker-down:
	@echo "🛑 Parando containers Docker..."
	@docker compose down

## docker-logs: Mostra logs dos containers
docker-logs:
	@docker compose logs -f

