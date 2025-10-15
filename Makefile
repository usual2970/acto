GO=go

.PHONY: run test lint swag mocks compose-up compose-down

run:
	$(GO) run ./app

test:
	$(GO) test ./...

lint:
	golangci-lint run

swag:
	swag init -g app/main.go -o ./docs

mocks:
	mockery

compose-up:
	docker compose up -d

compose-down:
	docker compose down
