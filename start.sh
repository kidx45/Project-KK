#!/bin/sh
set -e

echo "starting migration of database"
./migrate -path ./migration -database "$DB_URL" -verbose up

echo "starting the application"
exec "$@"