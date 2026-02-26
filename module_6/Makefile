.PHONY: test test-unit test-integration test-e2e test-all swagger

swagger:
	@echo "Generating Swagger docs..."
	~/go/bin/swag init -g cmd/server/main.go

test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/services/... ./cmd/server/handlers/... ./cmd/server/middlewares/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-e2e:
	@echo "Running E2E tests..."
	E2E_TEST=true go test -v ./tests/e2e/...

test-all:
	@echo "Running all tests..."
	go test -v ./...

test:
	@make test-unit
