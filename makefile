postgres: 
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb -U postgres -O postgres weight_tracker

dropdb:
	docker exec -it postgres12 dropdb -U postgres weight_tracker

migrateup:
	migrate -path pkg/repository/migrations -database "postgresql://postgres:postgres@localhost:5432/weight_tracker?sslmode=disable" -verbose up

migrateup1:
	migrate -path pkg/repository/migrations -database "postgresql://postgres:postgres@localhost:5432/weight_tracker?sslmode=disable" -verbose up 1

migratedown:
	migrate -path pkg/repository/migrations -database "postgresql://postgres:postgres@localhost:5432/weight_tracker?sslmode=disable" -verbose down

migratedown1:
	migrate -path pkg/repository/migrations -database "postgresql://postgres:postgres@localhost:5432/weight_tracker?sslmode=disable" -verbose down 1

test: 
	go test ./...
	

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 test