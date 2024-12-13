# name app
GOOSE_DBSTRING ?= "root:root1234@tcp(127.0.0.1:30306)/"
GOOSE_MIGRATION_DIR ?= sql/schema
GOOSE_DRIVER ?= mysql

APP_NAME := server
docker_build:
	docker-compose up -d --build
	docker-cmpose ps
docker_stop:
	docker-compose down
run:
	docker compose up -d && go run ./cmd/${APP_NAME}/user
dev:
	go run ./cmd/${APP_NAME}/$(name)
consumer:
	go run ./internal/consumer/main/
docker_up:
	docker compose up -d
up_by_one:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING)$(db_name) goose -dir=$(GOOSE_MIGRATION_DIR)/$(dir) up-by-one
create_migration:
	@goose -dir=$(GOOSE_MIGRATION_DIR)/$(dir) create $(name) sql
upse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING)$(db_name) goose -dir=$(GOOSE_MIGRATION_DIR)/$(dir) up
upto:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING)$(db_name) goose -dir=$(GOOSE_MIGRATION_DIR)/$(dir) up-to $(name)
downse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING)$(db_name) goose -dir=$(GOOSE_MIGRATION_DIR)/$(dir) down
resetse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING)$(db_name) goose -dir=$(GOOSE_MIGRATION_DIR)/$(dir) reset
sqlgen:
	sqlc generate
swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs
kafka_sync:
	curl -i -X POST -H "Accept: application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d @$(json_name)
.PHONY: dev downse upse resetse docker_build docker_stop docker_up swag upto