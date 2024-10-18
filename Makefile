.PHONY: postgres createdb dropdb migrate migrateup migratedown sqlc test migrategithubup migrategithub server mock

# Inicia los servicios de Docker
postgres:
	docker-compose up -d

# Crea la base de datos
createdb:
	docker exec -it postgresql-database createdb --username=alumno --owner=alumno simple_bank

# Elimina la base de datos
dropdb:
	docker exec -it postgresql-database dropdb --username=alumno simple_bank

# Ejecuta las migraciones
migrate:
	docker run --rm --network host -v ./db/migrations:/migrations migrate/migrate -path /migrations -database postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable up

migrategithub:
	migrate -path db/migrations -database "postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrategithubup:
	migrate -path db/migrations -database "postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup:
	docker run --rm --network host -v ./db/migrations:/migrations migrate/migrate -path /migrations -database postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable up

migratedown:
	docker run --rm --network host -v ./db/migrations:/migrations migrate/migrate -path /migrations -database postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable down -all

sqlc:
	sqlc generate

test:
	go test -v -cover ./...	

server:
	go run main.go

mock:
	mockgen -destination=db/mock/store.go -package=mockdb go_banking_system/db/sqlc Store