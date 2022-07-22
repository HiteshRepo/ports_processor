FROM golang:1.18-alpine AS builder

# Install `git` binary
RUN apk update
RUN apk add git build-base

## Install go tools
RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download
COPY .. .
RUN task ports_processor:compile

FROM alpine:3.14.2

## Install `wait` script
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

COPY --from=builder ./service/build/out/ports_processor ports_processor
COPY --from=builder ./service/configs/default.yaml default.yaml
COPY --from=builder ./service/internal/app/db/migrations/scripts scripts
COPY --from=builder ./service/data/ports.json ports.json

ENTRYPOINT ["./ports_processor", "-configFile", "default.yaml"]