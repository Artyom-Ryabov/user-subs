DB_CONNECTION="postgres://postgres:admin@localhost:5432/UserSubs"
MIGRATIONS=./db/migrations

run:
	go run .

sqlc:
	sqlc generate -f ./db/queries/sqlc.yaml

goose_up: 
	goose postgres ${DB_CONNECTION} up -dir ${MIGRATIONS}

goose_down: 
	goose postgres ${DB_CONNECTION} down -dir ${MIGRATIONS}
