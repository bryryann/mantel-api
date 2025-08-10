include .envrc

.PHONY: build
build:
	@echo 'Building binary into bin/...'
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api

.PHONY: run/api
run/api:
	go run ./cmd/api

.PHONY: db/migration
db/migration:
	@echo 'Creating migration: ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/up
db/up:
	@echo 'Setting up migrations...'
	migrate -path ./migrations -database ${MANTEL_DB_DSN} up

.PHONY: db/down
db/down:
	@echo 'Setting up migrations...'
	migrate -path ./migrations -database ${MANTEL_DB_DSN} down

.PHONY: test
test:
	go test -v -cover ./cmd/api/...

.PHONY: test-unit
test-unit:
	go test -v -short ./cmd/api/...

.PHONY: test-integration
test-integration:
	go test -v -run Integration ./cmd/api/...

.PHONY: bench
bench:
	go test -v -bench=. ./cmd/api/benchmark/...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./cmd/api/...
	go tool cover -html=coverage.out

