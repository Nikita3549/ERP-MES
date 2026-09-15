# erp-mes

## Run

Create `.env` (all variables are required):

```
DSN=postgres://user:pass@localhost:5432/erp-mes?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=pass
DB_NAME=erp-mes

HTTP_PORT=8080
HTTP_READ_HEADER_TIMEOUT=5s
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=20s
HTTP_IDLE_TIMEOUT=60s
HTTP_SHUTDOWN_TIMEOUT=20s
```

Docker (runs migrations, then the API):

```sh
docker compose --profile api up --build
```

Local:

```sh
docker compose up -d db
go build -o bin/api ./cmd && ./bin/api
```

Health check: `curl localhost:8080/health`

## Graceful shutdown

On SIGINT/SIGTERM the server drains in-flight requests for up to
`HTTP_SHUTDOWN_TIMEOUT`, closes the DB and exits `0` (`1` on failure).
Keep `HTTP_SHUTDOWN_TIMEOUT >= HTTP_WRITE_TIMEOUT` and below `stop_grace_period` (25s).

Check (expect `Server stopped`, exit 0; do not use `go run`, it swallows SIGTERM):

```sh
./bin/api & sleep 1
kill -TERM $!; wait $!; echo $?
```
