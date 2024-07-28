build:
	@go build -o bin/rental_listing

run: build
	 @./bin/rental_listing

postgres:
	@docker run --name postgresv3 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	@docker exec -it postgresv3 createdb --username=root --owner=root rental_listing

dropdb:
	@docker exec -it postgresv3 dropdb rental_listing

DB_URL=postgresql://root:secret@localhost:5432/rental_listing?sslmode=disable

migrateup:
	@migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	@migrate -path db/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	@migrate -path db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	@migrate -path db/migrations -database "$(DB_URL)" -verbose down 1 

migratedirty:
	@migrate -path db/migrations -database "$(DB_URL)" force 1


test:
	@go test -v -cover -short ./...

sqlc:
	@sqlc generate

new_migration:
	@migrate create -ext sql -dir db/migrations -seq $(name)


redis:
	@docker run --name redis -p 6379:6379 -d redis:7-alpine

# command line tool to inspect the state of queues and tasks.
asynq:
	@asynq dash
