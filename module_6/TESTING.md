# Testing Guide

## Test Coverage

### Unit Tests
- `internal/services/auth_test.go` - AuthService tests
- `internal/services/chat_test.go` - ChatService tests
- `cmd/server/handlers/auth_test.go` - Auth handler tests
- `cmd/server/handlers/chat_test.go` - Chat handler tests
- `cmd/server/middlewares/auth_test.go` - Auth middleware tests

### Integration Tests
- `tests/integration/integration_test.go` - Full flow tests

### E2E Tests
- `tests/e2e/e2e_test.go` - End-to-end tests with running server

## Running Tests

### Unit Tests Only
```bash
make test-unit
# or
go test -v ./internal/services/... ./cmd/server/handlers/... ./cmd/server/middlewares/...
```

### Integration Tests
```bash
make test-integration
# or
go test -v ./tests/integration/...
```

### E2E Tests
Start the server first:
```bash
docker-compose up -d
go run cmd/server/main.go
```

Then run E2E tests:
```bash
make test-e2e
# or
E2E_TEST=true go test -v ./tests/e2e/...
```

### All Tests
```bash
make test-all
# or
go test -v ./...
```

## Requirements
- MongoDB running on `localhost:27017` for unit/integration tests
- Server running on `localhost:8080` for E2E tests
