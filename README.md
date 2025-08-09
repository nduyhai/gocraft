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
cleanctl new myapp -m github.com/you/myapp -t basic -o ./myapp
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