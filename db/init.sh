#!/bin/bash
set -e

echo "Running database migrations..."

for f in /docker-entrypoint-initdb.d/migrations/*.up.sql; do
  echo "Applying: $(basename "$f")"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done

echo "Running seed data..."

for f in /docker-entrypoint-initdb.d/seeds/*.sql; do
  echo "Seeding: $(basename "$f")"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done

echo "Database initialisation complete."
