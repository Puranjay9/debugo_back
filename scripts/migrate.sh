#!/usr/bin/env bash
set -e

NAME=$1

if [ -z "$NAME" ]; then
  echo "❌ Migration name required"
  echo "Usage: ./scripts/migrate.sh add_users_table"
  exit 1
fi

echo "📦 Creating migration: $NAME"

migrate create \
  -ext sql \
  -dir migrations \
  -seq "$NAME"

echo "✏️  Edit the migration files now"
read -p "Press ENTER when ready to apply migrations..."

echo "🚀 Applying migrations via Docker"
docker compose run --rm migrate

echo "✅ Migration applied successfully"
