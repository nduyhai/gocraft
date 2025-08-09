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
