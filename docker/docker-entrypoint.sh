#!/bin/sh
set -e

# Миграция
migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./migrations up

exec "$@" 
