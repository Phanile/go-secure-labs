include .env
export

APP_NAME := go-secure-labs
PROJECT_ROOT := $(shell pwd)
CMD_FILE := ${PROJECT_ROOT}/cmd/main.go
LOGGER_FOLDER := ${PROJECT_ROOT}/out/logs

.PHONY: build run test

run:
	go run $(CMD_FILE)

env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

env-clean:
	@read -p "Are you sure you want to delete the database data? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down postgres && \
		rm -rf out/pgdata && \
		echo "Database data deleted."; \
	else \
		echo "Database data deletion canceled."; \
	fi

migrate:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Please provide a seq for the migration. Example: make migrate seq=init"; \
		exit 1; \
	fi
	@mkdir -p migrations
	@echo "Creating database migrations..."
	@docker compose run --rm -q \
	--user "$$(id -u):$$(id -g)" \
	postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"
	
migrate-up:
	@echo "Applying database migrations..."
	@make migrate-action action=up

migrate-down:
	@echo "Reverting database migrations..."
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Please provide an action for the migration. Example: make migrate-action action=up"; \
		exit 1; \
	fi
	@HOST_UID=$$(id -u) HOST_GID=$$(id -g) docker compose run --rm postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
	"$(action)"

deploy:
	@docker compose up -d --build go-secure-labs