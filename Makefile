postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mahanta -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root e_commerce

dropdb: 
	docker exec -it postgres17 dropdb e_commerce

migrateup:
	migrate -path db/migration -database "postgresql://root:mahanta@localhost:5432/e_commerce?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:mahanta@localhost:5432/e_commerce?sslmode=disable" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

run:
	go run main.go

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Yelsnik/e-commerce-api/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown