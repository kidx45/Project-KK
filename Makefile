postgres:
	docker run --name postgres18 --network k_network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine 
createdb:
	docker exec -it postgres18 createdb --username=root --owner=root project_test_kk
dropdb:
	docker exec -it postgres18 dropdb project_test_kk
migrateup:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5434/project_kk?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration -database "postgresql://root:secret@localhost:5434/project_kk?sslmode=disable" -verbose down
migrateupmysql:
	migrate -path ./db/migration-mysql -database "mysql://root:secret@tcp(localhost:3307)/project_kk?tls=false" -verbose up
migratedownmysql:
	migrate -path ./db/migration-mysql -database "mysql://root:secret@tcp(localhost:3307)/project_kk?tls=false" -verbose down
migrateCreate:
	migrate create -ext sql -dir ./db/migration -seq schema_lifeline
migrateCreatemysql:
	migrate create -ext sql -dir ./db/migration-mysql -seq schema_lifeline_mysql
sqlc:
	sqlc generate
test:
	go test -v -cover -coverpkg=./... ./...
mockgen:
	mockgen -package=mockdb -destination=./db/mockdb/store.go github.com/kidx45/Project-KK/db/sqlc Store
dockerize:
	sudo docker build -t project_kk:1.0 .
dockerize-run:
	sudo docker run --name kk --network k_network -p 8080:8080 -e GIN_MODE=release -e DB_URL="postgresql://root:secret@psql_testing:5432/project_kk?sslmode=disable" project_kk:1.0
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb \
    --grpc-gateway_opt=paths=source_relative \
    proto/*.proto
.PHONY: postgres createdb dropdb migrateup migratedown migrateCreate sqlc proto
