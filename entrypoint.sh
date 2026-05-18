#!/bin/sh
set -e

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=postgres port=5432 user=postgres password=postgres dbname=gogen sslmode=disable"
export GOOSE_MIGRATION_DIR=./sql/migrations

echo "Running migrations..."
goose up

echo "Starting server..."
exec ./server