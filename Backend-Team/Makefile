postgres:
	docker run --name postgres18 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine 
createdb:
	docker exec -it postgres18 createdb --username=root --owner=root Project_KK
dropdb:
	docker exec -it postgres18 dropdb Project_KK
migrateup:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5433/Project_KK?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5433/Project_KK?sslmode=disable" -verbose down
migrateCreate:
	migrate create -ext sql -dir ./db/migration -seq schema_lifeline
sqlc:
	sqlc generate
.PHONY: postgres createdb dropdb migrateup migratedown migrateCreate sqlc