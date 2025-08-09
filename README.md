# cleanctl

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/ynduyhai/go-clean-arch-starter)](LICENSE)

A CLI to generate Go projects with a clean architecture layout from embedded templates.

## Install

```bash
go install github.com/nduyhai/go-clean-arch-starter/cmd/cleanctl@latest
```

## Usage

```bash
cleanctl new myapp -m github.com/you/myapp
```

This generates:
- go.mod (module set to github.com/you/myapp)
- cmd/api/main.go
- README.md

Then optionally initialize git and tidy dependencies automatically.

## Structure

See internal directory for core, adapters, and platform layers. Templates are embedded in internal/adapters/outbound/templates/embed_repo/templates.


```
<your-app>/
├─ cmd/
│  └─ api/                    # Composition root (wire deps, start servers/cron/CMDs)
│     └─ main.go
├─ internal/
│  ├─ core/                   # Enterprise/business rules (NO framework deps)
│  │  ├─ entity/              # Domain entities
│  │  │  └─ user.go
│  │  ├─ ports/               # Inbound/Outbound interfaces
│  │  │  └─ user_repository.go
│  │  └─ usecase/             # Application services (orchestrate domain logic)
│  │     └─ user_service.go
│  ├─ adapters/               # Implements ports (drive & driven adapters)
│  │  ├─ repository/
│  │  │  ├─ memory/           # In-memory (default)
│  │  │  │  └─ user_repo.go
│  │  │  └─ postgres/         # Optional: build-tagged or generated when selected
│  │  │     └─ user_repo.go
│  │  ├─ transport/
│  │  │  ├─ http/             # HTTP handlers, router, DTOs, validation
│  │  │  │  ├─ router.go
│  │  │  │  └─ user_handler.go
│  │  │  └─ grpc/             # Optional: gRPC server & proto adapters
│  │  └─ cache/
│  │     └─ redis/            # Optional: Redis adapter (e.g., for read-cache)
│  └─ platform/               # Cross-cutting utilities (infra-agnostic)
│     ├─ config/              # Env/config loading
│     ├─ logger/              # slog setup, contexts, fields
│     └─ id/                  # UUID/KSUID generator (so usecase doesn’t import libs)
├─ pkg/                       # (Optional) public libs usable by other modules
├─ migrations/                # DB migrations (goose, atlas, flyway…)
├─ api/                       # OpenAPI/gRPC IDL & generated clients/servers
│  ├─ openapi/
│  └─ proto/
├─ test/                      # Integration/end-to-end test helpers
├─ .env.example
├─ docker-compose.yml
├─ Makefile
├─ go.mod
└─ README.md

```

## Module

### Core
| Module Name              | Purpose                                 | Requires        | Conflicts |
| ------------------------ | --------------------------------------- | --------------- | --------- |
| `platform:base`          | Fx + Viper config, logger, DI root      | –               | –         |
| `platform:logger`        | Structured logging (slog or zap)        | `platform:base` | –         |
| `platform:testing`       | Unit test helpers, testcontainers setup | `platform:base` | –         |
| `platform:observability` | OpenTelemetry tracing, metrics, pprof   | `platform:base` | –         |


### Transports

| Module Name   | Purpose                                 | Requires        | Conflicts               |
| ------------- | --------------------------------------- | --------------- | ----------------------- |
| `http:gin`    | HTTP server via Gin, DI lifecycle       | `platform:base` | `http:chi`, `http:echo` |
| `http:chi`    | HTTP server via Chi                     | `platform:base` | `http:gin`, `http:echo` |
| `grpc:server` | gRPC server transport                   | `platform:base` | –                       |
| `grpc:client` | gRPC client support                     | `platform:base` | –                       |
| `rest:client` | HTTP client with retry, circuit breaker | `platform:base` | –                       |


### Database

| Module Name    | Purpose                           | Requires        | Conflicts                 |
| -------------- | --------------------------------- | --------------- | ------------------------- |
| `db:postgres`  | Postgres adapter (pgx/sqlc)       | `platform:base` | `db:mysql`, `db:gorm`     |
| `db:mysql`     | MySQL adapter (sqlc/mysql driver) | `platform:base` | `db:postgres`, `db:gorm`  |
| `db:gorm`      | ORM with GORM                     | `platform:base` | `db:postgres`, `db:mysql` |
| `db:sqlite`    | SQLite for local dev/tests        | `platform:base` | –                         |
| `db:migration` | Goose or Atlas migrations         | `platform:base` | –                         |


### Caching

| Module Name      | Purpose                    | Requires        | Conflicts |
| ---------------- | -------------------------- | --------------- | --------- |
| `cache:redis`    | Redis connection + helper  | `platform:base` | –         |
| `queue:kafka`    | Kafka producer/consumer    | `platform:base` | –         |
| `queue:rabbitmq` | RabbitMQ producer/consumer | `platform:base` | –         |
| `pubsub:nats`    | NATS JetStream setup       | `platform:base` | –         |


### Utilities

| Module Name       | Purpose                          | Requires                 | Conflicts |
| ----------------- | -------------------------------- | ------------------------ | --------- |
| `feature:i18n`    | i18n JSON file loader            | `platform:base`          | –         |
| `feature:auth`    | JWT auth middleware              | `http:gin` or `http:chi` | –         |
| `feature:health`  | `/health` and `/ready` endpoints | `http:*`                 | –         |
| `feature:metrics` | Prometheus metrics endpoint      | `http:*`                 | –         |
| `feature:swagger` | Swagger/OpenAPI docs generation  | `http:*`                 | –         |
