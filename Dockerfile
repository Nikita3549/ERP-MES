FROM golang:1.26.4 as build
WORKDIR /opt/api

COPY go.mod go.sum ./

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN go build ./cmd/main.go

FROM node:22-slim
WORKDIR /opt/api

COPY --chown=node:node ./migrations /opt/api/migrations
COPY --from=build --chown=node:node /opt/api/main /opt/api/main
COPY --from=build /go/bin/migrate /usr/local/bin/migrate

USER node

CMD . ./.env \
  && migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up \
  && exec ./main

