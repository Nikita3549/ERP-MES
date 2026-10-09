# erp-mes

## Run

Docker (db → migrations → api):

    docker compose --profile api up --build

Local:

    docker compose up -d db
    go build -o bin/api ./cmd && ./bin/api

## Config

Env vars or `.env` (see `.env.example`). All optional:

| Var | Default |
|---|---|
| DB_HOST | localhost |
| DB_PORT | 5432 |
| DB_USER | postgres |
| DB_PASSWORD | postgres |
| DB_NAME | erp-mes |
| HTTP_PORT | 8080 |
| HTTP_READ_HEADER_TIMEOUT | 5s |
| HTTP_READ_TIMEOUT | 15s |
| HTTP_WRITE_TIMEOUT | 20s |
| HTTP_IDLE_TIMEOUT | 60s |
| HTTP_SHUTDOWN_TIMEOUT | 20s |

On SIGINT/SIGTERM the server drains requests for up to `HTTP_SHUTDOWN_TIMEOUT`, then closes the DB.
Keep it `>= HTTP_WRITE_TIMEOUT` and below compose `stop_grace_period` (25s).

## Endpoints

- `GET /health` — liveness, `{"status":"ok"}`
- `GET /ready` — readiness, pings DB; `503 {"status":"fail"}` on failure

## Migrations

`migrations/NNNNNN_<name>.{up,down}.sql`, applied by `migrate/migrate` before the API starts.
