#!/bin/sh
set -e

echo "⏳ Waiting for PostgreSQL to start..."

# Wait until Postgres port is reachable
while ! nc -z postgres 5432; do
  sleep 1
done

echo "🐘 PostgreSQL is up - executing database migrations..."
./migrate

echo "🚀 Starting Go Fiber Backend..."
exec ./api