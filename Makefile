.PHONY: run test lint up down generate migrate

run:
	go run ./cmd/main.go

test:
	go test ./... -v

test-short:
	go test ./... -short

lint:
	golangci-lint run ./...

generate:
	buf generate

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f app

migrate:
	goose -dir migrations postgres "host=localhost port=5432 user=marketplace password=marketplace dbname=marketplace sslmode=disable" up
