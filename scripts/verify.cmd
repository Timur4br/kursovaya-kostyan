@echo off
setlocal
cd /d "%~dp0.."
if not exist ".env" (
  echo Copy .env.sample to .env and fill all four values first.
  exit /b 1
)
if not exist ".local" mkdir ".local"
set "PERSISTENCE_PROBE="
echo Checking Compose configuration...
docker compose config --quiet > ".local\verify.log" 2>&1
if errorlevel 1 goto failed
echo Building and starting PostgreSQL, migrations and API. This can take several minutes...
docker compose up --build -d --wait --wait-timeout 120 >> ".local\verify.log" 2>&1
if errorlevel 1 goto failed
echo Running Go tests and PostgreSQL constraint checks...
docker compose --profile test run --build --rm tests >> ".local\verify.log" 2>&1
if errorlevel 1 goto failed
set "PERSISTENCE_PROBE=local-%RANDOM%-%RANDOM%"
docker compose --profile test run --rm tests go test -count=1 -run "^TestPersistenceWrite$" ./services/api_gateway/internal/database >> ".local\verify.log" 2>&1
if errorlevel 1 goto failed
echo Restarting this project's database and checking stored data...
docker compose restart db >> ".local\verify.log" 2>&1
if errorlevel 1 goto failed
docker compose up -d --wait --wait-timeout 120 >> ".local\verify.log" 2>&1
if errorlevel 1 goto failed
docker compose --profile test run --rm tests go test -count=1 -run "^TestPersistenceRead$" ./services/api_gateway/internal/database >> ".local\verify.log" 2>&1
if errorlevel 1 goto failed
echo PASS - Compose startup, migrations, schema constraints and persistence verified.
echo Full log: .local\verify.log
docker compose ps
exit /b 0
:failed
echo FAIL - verification stopped. Full log follows:
type ".local\verify.log"
exit /b 1
