#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
unset PERSISTENCE_PROBE
docker compose config --quiet
docker compose up --build -d --wait --wait-timeout 120
docker compose --profile test run --build --rm tests
PERSISTENCE_PROBE="ci-$(date +%s)-$$"
export PERSISTENCE_PROBE
docker compose --profile test run --rm tests go test -count=1 -run '^TestPersistenceWrite$' ./services/api_gateway/internal/database
docker compose restart db
docker compose up -d --wait --wait-timeout 120
docker compose --profile test run --rm tests go test -count=1 -run '^TestPersistenceRead$' ./services/api_gateway/internal/database
docker compose ps
echo 'PASS - Compose startup, migrations, schema constraints and persistence verified.'
