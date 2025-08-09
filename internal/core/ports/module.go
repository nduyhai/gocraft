package ports

type Ctx interface {
}

type Module interface {
	Name() string  // machine-friendly, unique, e.g. "db:postgres"
	Label() string // user-friendly, e.g. "PostgreSQL Adapter (pgx/sqlc)"
	Version() string
	Summary() string // short one-line description
	Tags() []string

	Requires() []string
	Conflicts() []string
	Applies(ctx Ctx) bool
	Apply(ctx Ctx) error
}
