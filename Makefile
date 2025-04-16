DB_URL=postgres://root:secret@localhost:5432/p2platform?sslmode=disable
migrateinstall:
	$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migrateuplast:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedownlast:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
postgres:
	docker run --name postgres17 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:17-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root p2platform
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go p2platform/db/sqlc Store

.PHONY: migrateinstall migrateup new_migration postgres createdb migratedown sqlc test server mock