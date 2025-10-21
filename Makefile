-include .env
DB_CONNECTION="postgres://${DB_USER}:${DB_PSWD}@${DB_HOST}:${DB_PORT}/UserSubs"
MIGRATIONS=./migrations

sqlc:
	sqlc generate -f ./db/queries/sqlc.yaml

goose_up: 
	goose postgres ${DB_CONNECTION} up -dir ${MIGRATIONS}

goose_down: 
	goose postgres ${DB_CONNECTION} down -dir ${MIGRATIONS}

swagger:
	swag i -g ./cmd/main.go

run: goose_down goose_up sqlc swagger
	go run ./cmd/main.go

build: goose_down goose_up sqlc swagger
	go build -C ./cmd -o main

