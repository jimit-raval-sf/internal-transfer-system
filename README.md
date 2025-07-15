# Internal Transfer System

A Go-based internal transfer system built with Gin and GORM, backed by PostgreSQL.

## Features

- Account creation with validation
- Account balance queries
- Internal money transfers with atomicity
- Clean architecture with proper error handling
- Graceful shutdown
- Database transactions for consistency

## API Endpoints

### 1. Create Account
**POST** `/api/v1/accounts`

```json
{
  "account_id": 123,
  "initial_balance": "100.23344"
}
```

**Responses:**
- `201`: Account created successfully
- `400`: Invalid input (account_id must be positive, balance must be non-negative with max 5 decimal places)
- `409`: Account already exists

### 2. Get Account
**GET** `/api/v1/accounts/{account_id}`

**Response:**
```json
{
  "account_id": 123,
  "balance": "100.23344"
}
```

**Responses:**
- `200`: Account found
- `404`: Account not found

### 3. Create Transaction
**POST** `/api/v1/transactions`

```json
{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "100.12345"
}
```

**Responses:**
- `200`: Transaction successful
- `400`: Invalid input (amount must be > 0, no self-transfer)
- `404`: Account(s) not found
- `422`: Insufficient balance

## Setup

1. Start PostgreSQL:
```bash
docker-compose up -d
```

2. Configure environment variables (choose one method):

   **Option A: Using .env file (recommended)**
   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

   **Option B: Using environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=internal_transfer
   export DB_PORT=5432
   export PORT=8080
   ```

3. Run the application:
```bash
go mod tidy
go run main.go
```

## Testing

Example requests:

```bash
# Create accounts
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "initial_balance": "1000.00"}'

curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{"account_id": 2, "initial_balance": "500.00"}'

# Get account
curl http://localhost:8080/api/v1/accounts/1

# Transfer money
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{"source_account_id": 1, "destination_account_id": 2, "amount": "100.00"}'
```

## Architecture

- **Clean Architecture**: Separated into handlers, services, and repositories
- **Database**: PostgreSQL with GORM ORM
- **Web Framework**: Gin for HTTP routing
- **Decimal Handling**: shopspring/decimal for precise monetary calculations
- **Concurrency**: Database-level locking for transaction safety
- **Graceful Shutdown**: Proper server shutdown handling