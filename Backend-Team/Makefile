postgres:
	docker run --name postgres18 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine 
createdb:
	docker exec -it postgres18 createdb --username=root --owner=root project_test_kk
dropdb:
	docker exec -it postgres18 dropdb project_test_kk
migrateup:
	migrate -path ./db/migration -database "postgres://root:secret@localhost:5433/project_test_kk?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration -database "postgres://root:secret@localhost:5433/project_test_kk?sslmode=disable" -verbose down
migrateupmysql:
	migrate -path ./db/migration-mysql -database "mysql://root:secret@tcp(localhost:3307)/project_test_kk?tls=false" -verbose up
migratedownmysql:
	migrate -path ./db/migration-mysql -database "mysql://root:secret@tcp(localhost:3307)/project_test_kk?tls=false" -verbose down
migrateCreate:
	migrate create -ext sql -dir ./db/migration -seq schema_lifeline
migrateCreatemysql:
	migrate create -ext sql -dir ./db/migration-mysql -seq schema_lifeline_mysql
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown migrateCreate sqlc