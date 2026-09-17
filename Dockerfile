FROM golang:1.26.4 as build
WORKDIR /opt/api

COPY go.mod go.sum ./

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 go build -o main ./cmd

FROM alpine:latest
WORKDIR /opt/api

RUN adduser -D -H app

COPY --from=build --chown=app:app /opt/api/main /opt/api/main

USER app
ENTRYPOINT [ "./main" ]
