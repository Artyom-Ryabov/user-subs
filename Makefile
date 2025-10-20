DB_CONNECTION="postgres://postgres:admin@localhost:5432/UserSubs"
MIGRATIONS=./migrations

run:
	go run ./cmd/main.go

sqlc:
	sqlc generate -f ./db/queries/sqlc.yaml

goose_up: 
	goose postgres ${DB_CONNECTION} up -dir ${MIGRATIONS}

goose_down: 
	goose postgres ${DB_CONNECTION} down -dir ${MIGRATIONS}

swagger:
	swag i -g ./cmd/main.go
