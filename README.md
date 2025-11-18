# Document Validation API

API REST para validação e gerenciamento de CPF/CNPJ desenvolvida em Golang com Fiber, GORM e PostgreSQL.

## Funcionalidades

- Validação de CPF/CNPJ com verificação de dígitos verificadores
- CRUD completo de documentos (CPF/CNPJ)
- Filtros e ordenação na listagem
- Sistema de blocklist para marcar documentos
- Autenticação JWT
- Endpoint de status com uptime e contador de requisições
- Testes unitários

## Tecnologias

- **Golang 1.21**
- **Fiber** - Framework web
- **GORM** - ORM para banco de dados
- **PostgreSQL** - Banco de dados
- **JWT** - Autenticação
- **Zerolog** - Logging
- **Docker & Docker Compose** - Conteinerização

## Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.21+ (para desenvolvimento local)

## Instalação e Execução

### Usando Docker Compose (Recomendado)

1. Clone o repositório:
```bash
git clone <repository-url>
cd golang-validation-documents
```

2. Execute o projeto com Docker Compose:
```bash
docker-compose up --build
```

A API estará disponível em `http://localhost:8080`

### Desenvolvimento Local

1. Instale as dependências:
```bash
go mod download
```

2. Configure as variáveis de ambiente (opcional, valores padrão já configurados):
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=validation_user
export DB_PASSWORD=validation_pass
export DB_NAME=validation_db
export DB_SSLMODE=disable
export JWT_SECRET=your-secret-key-change-in-production
export PORT=8080
```

3. Execute o PostgreSQL localmente ou use Docker:
```bash
docker run -d \
  --name validation_postgres \
  -e POSTGRES_USER=validation_user \
  -e POSTGRES_PASSWORD=validation_pass \
  -e POSTGRES_DB=validation_db \
  -p 5432:5432 \
  postgres:15-alpine
```

4. Execute a aplicação:
```bash
go run cmd/api/main.go
```

## Endpoints da API

### Autenticação

#### POST /login
Autentica e retorna um token JWT.

**Request:**
```json
{
  "username": "admin",
  "password": "admin"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Status

#### GET /status
Retorna informações sobre o status do servidor.

**Response:**
```json
{
  "uptime": "2h30m15s",
  "uptime_seconds": 9015,
  "request_count": 42,
  "status": "ok"
}
```

### Documentos

Todos os endpoints de documentos requerem autenticação via header:
```
Authorization: Bearer <token>
```

#### POST /api/v1/documents
Cria um novo documento (CPF/CNPJ).

**Request:**
```json
{
  "number": "11144477735"
}
```

**Response:**
```json
{
  "id": "uuid",
  "number": "11144477735",
  "type": "CPF",
  "blocked": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### GET /api/v1/documents
Lista documentos com filtros e paginação.

**Query Parameters:**
- `number` (opcional): Filtrar por número
- `type` (opcional): Filtrar por tipo (CPF ou CNPJ)
- `blocked` (opcional): Filtrar por status de bloqueio (true/false)
- `order_by` (opcional): Ordenação (ex: "created_at DESC", "number ASC")
- `page` (opcional): Número da página (padrão: 1)
- `page_size` (opcional): Itens por página (padrão: 10, máximo: 100)

**Exemplo:**
```
GET /api/v1/documents?type=CPF&blocked=false&order_by=created_at DESC&page=1&page_size=10
```

**Response:**
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

#### GET /api/v1/documents/search
Busca um documento pelo número.

**Query Parameters:**
- `number` (obrigatório): Número do documento

**Exemplo:**
```
GET /api/v1/documents/search?number=11144477735
```

#### GET /api/v1/documents/:id
Busca um documento pelo ID.

**Response:**
```json
{
  "id": "uuid",
  "number": "11144477735",
  "type": "CPF",
  "blocked": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /api/v1/documents/:id
Atualiza um documento (principalmente para blocklist).

**Request:**
```json
{
  "blocked": true
}
```

**Response:**
```json
{
  "id": "uuid",
  "number": "11144477735",
  "type": "CPF",
  "blocked": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:01Z"
}
```

#### DELETE /api/v1/documents/:id
Remove um documento (soft delete).

**Response:** 204 No Content

## Testes

Execute os testes unitários de forma simples:

### Usando Make (Recomendado - mais simples)

```bash
make test          # Executa todos os testes
make coverage      # Executa testes com cobertura
make test-verbose  # Executa testes com output detalhado
make test-utils    # Executa apenas testes do pacote utils
make test-service  # Executa apenas testes do pacote service
make help          # Mostra todos os comandos disponíveis
```

### Comandos Docker diretos (alternativa)

Se preferir não usar Make, você pode rodar diretamente:

```bash
# Executar todos os testes
docker run --rm -v "$(pwd)":/app -w /app golang:1.23-alpine go test ./...

# Executar testes com cobertura
docker run --rm -v "$(pwd)":/app -w /app golang:1.23-alpine go test -cover ./...

# Executar testes com output detalhado
docker run --rm -v "$(pwd)":/app -w /app golang:1.23-alpine go test -v ./...
```

### Desenvolvimento Local (com Go instalado)

Se você tiver Go instalado na sua máquina:

```bash
go test ./...
go test -cover ./...
```

## Estrutura do Projeto

```
.
├── cmd/
│   └── api/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go            # Configurações da aplicação
│   ├── database/
│   │   └── database.go          # Conexão com banco de dados
│   ├── domain/
│   │   └── document.go          # Entidade Document
│   ├── handler/
│   │   ├── auth_handler.go      # Handlers de autenticação
│   │   ├── document_handler.go  # Handlers de documentos
│   │   └── status_handler.go    # Handler de status
│   ├── middleware/
│   │   ├── auth.go              # Middleware de autenticação JWT
│   │   └── security.go          # Middlewares de segurança
│   ├── repository/
│   │   └── document_repository.go # Camada de acesso a dados
│   ├── service/
│   │   ├── document_service.go  # Lógica de negócio
│   │   └── document_service_test.go # Testes do serviço
│   └── utils/
│       ├── validator.go         # Validação de CPF/CNPJ
│       └── validator_test.go    # Testes de validação
├── migrations/
│   ├── 000001_create_documents.up.sql
│   └── 000001_create_documents.down.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## Validação de CPF/CNPJ

A validação implementa o algoritmo oficial de verificação de dígitos verificadores:

- **CPF**: Valida os dois dígitos verificadores
- **CNPJ**: Valida os dois dígitos verificadores com pesos específicos

Documentos com todos os dígitos iguais são considerados inválidos.

## Segurança

- Autenticação JWT obrigatória para endpoints protegidos
- Sanitização de inputs para prevenir XSS
- Validação de dados de entrada
- Soft delete para preservar histórico

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| DB_HOST | Host do PostgreSQL | localhost |
| DB_PORT | Porta do PostgreSQL | 5432 |
| DB_USER | Usuário do banco | validation_user |
| DB_PASSWORD | Senha do banco | validation_pass |
| DB_NAME | Nome do banco | validation_db |
| DB_SSLMODE | Modo SSL | disable |
| JWT_SECRET | Chave secreta JWT | your-secret-key-change-in-production |
| PORT | Porta da API | 8080 |

## Licença

Este projeto foi desenvolvido como teste técnico.
