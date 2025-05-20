include .envrc

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
