package gormmodule

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"gopkg.in/yaml.v3"
)

// Module implements ports.Module for adding a GORM setup with multi-driver support.
//
// Name:      db:gorm
// Requires:  platform:base
// Conflicts: db:postgres, db:mysql (since this provides DB access via GORM)
//
// This module adds:
// - internal/platform/db/gorm/ (Fx provider wiring *gorm.DB)
// - config defaults for gorm: driver + DSNs/paths for drivers
// - go.mod deps for gorm and drivers

type Module struct{}

func New() Module { return Module{} }

func (Module) Name() string    { return "db:gorm" }
func (Module) Label() string   { return "GORM (multi-driver: postgres/mysql/sqlite)" }
func (Module) Version() string { return "0.1.0" }
func (Module) Summary() string {
	return "Adds GORM with runtime-selectable driver via config/env (postgres/mysql/sqlite)"
}
func (Module) Tags() []string { return []string{"db", "gorm", "orm"} }

func (Module) Requires() []string  { return []string{"platform:base"} }
func (Module) Conflicts() []string { return []string{"db:postgres", "db:mysql"} }

func (Module) Applies(ctx ports.Ctx) bool { return true }

func (Module) Apply(ctx ports.Ctx) error {
	// Try to add required dependencies if a GoMod editor is available
	if gm := ctx.GoMod(); gm != nil {
		// Always add core gorm and infra libs
		_ = gm.Add("gorm.io/gorm", "v1.25.7-0.20240204074919-46816ad31dde")
		_ = gm.Add("github.com/spf13/viper", "v1.20.1")
		_ = gm.Add("go.uber.org/fx", "v1.24.0")

		// Add only the chosen SQL driver based on config/DSN
		drv := strings.ToLower(strings.TrimSpace(nestedString(ctx.Values(), []string{"gorm", "driver"})))
		if drv == "" {
			dsn := nestedString(ctx.Values(), []string{"gorm", "dsn"})
			drv = driverFromDSN(dsn)
		}
		if drv == "" {
			drv = "sqlite"
		}
		switch drv {
		case "postgres", "pg", "postgre", "postgresql":
			_ = gm.Add("gorm.io/driver/postgres", "v1.5.7")
		case "mysql":
			_ = gm.Add("gorm.io/driver/mysql", "v1.5.6")
		default: // sqlite
			_ = gm.Add("gorm.io/driver/sqlite", "v1.5.7")
		}
	}

	// Render templates from embedded FS
	sub, err := fs.Sub(TemplatesFS, "templates")
	if err != nil {
		return fmt.Errorf("sub fs: %w", err)
	}
	var tfiles []ports.TmplFile
	err = fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(sub, path)
		if err != nil {
			return err
		}
		tfiles = append(tfiles, ports.TmplFile{Path: path, Content: string(b)})
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk: %w", err)
	}
	tpl := ports.Template{Name: "db:gorm", Files: tfiles}
	files, err := ctx.Renderer().Render(tpl, ctx.Values())
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := ctx.FS().WriteAll(ctx.ProjectRoot(), files); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// Update DI root to include gorm module
	if ed := ctx.AdaptersModule(); ed != nil {
		vals := ctx.Values()
		modPath, _ := vals["Module"].(string)
		if modPath == "" {
			return fmt.Errorf("module path missing in context")
		}
		if err := ed.Ensure("gormdb", modPath+"/internal/platform/db/gorm", "gormdb.Module()"); err != nil {
			return fmt.Errorf("update di root: %w", err)
		}
	}

	// Respect --set gorm.driver=... by updating config/config.yml if provided.
	drv := nestedString(ctx.Values(), []string{"gorm", "driver"})
	if drv != "" {
		_ = ensureGormDriverInConfig(ctx.ProjectRoot(), drv)
	} else {
		// If no driver but DSN present, infer and set it in config
		dsn := nestedString(ctx.Values(), []string{"gorm", "dsn"})
		if dsn != "" {
			if inf := driverFromDSN(dsn); inf != "" {
				_ = ensureGormDriverInConfig(ctx.ProjectRoot(), inf)
			}
		}
	}
	return nil
}

// Defaults implements ports.Module.Defaults to provide default configuration.
func (Module) Defaults() map[string]any {
	// Load static defaults from embedded YAML if present
	if b, err := fs.ReadFile(TemplatesFS, "templates/config/defaults.yml.tmpl"); err == nil {
		var m map[string]any
		if yaml.Unmarshal(b, &m) == nil && m != nil {
			return m
		}
	}
	return map[string]any{
		"gorm": map[string]any{
			"driver":            "sqlite",
			"dsn":               "file:app.db?_pragma=busy_timeout=5000&_pragma=journal_mode=WAL",
			"log_level":         "warn",
			"max_open_conns":    25,
			"max_idle_conns":    5,
			"conn_max_lifetime": "30m",
		},
	}
}

func nestedString(m map[string]any, path []string) string {
	cur := any(m)
	for i, p := range path {
		mm, ok := cur.(map[string]any)
		if !ok {
			return ""
		}
		v, ok := mm[p]
		if !ok {
			return ""
		}
		if i == len(path)-1 {
			s, _ := v.(string)
			return s
		}
		cur = v
	}
	return ""
}

func ensureGormDriverInConfig(root, driver string) error {
	p := root + "/config/config.yml"
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	var current map[string]any
	if len(b) > 0 {
		_ = yaml.Unmarshal(b, &current)
	}
	if current == nil {
		current = make(map[string]any)
	}
	setPath(current, []string{"gorm", "driver"}, driver)
	out, err := yaml.Marshal(current)
	if err != nil {
		return err
	}
	return os.WriteFile(p, out, 0o644)
}

func setPath(m map[string]any, path []string, value any) {
	cur := m
	for i, p := range path {
		if i == len(path)-1 {
			cur[p] = value
			return
		}
		next, ok := cur[p]
		if !ok {
			nm := make(map[string]any)
			cur[p] = nm
			cur = nm
			continue
		}
		nm, ok := next.(map[string]any)
		if !ok {
			nm = make(map[string]any)
			cur[p] = nm
		}
		cur = nm
	}
}

// driverFromDSN attempts to infer the DB driver from a DSN string.
// Returns one of: "postgres", "mysql", "sqlite"; empty string if unknown.
func driverFromDSN(dsn string) string {
	dsn = strings.TrimSpace(strings.ToLower(dsn))
	if dsn == "" {
		return ""
	}
	// Postgres: URI schemes or keyword DSN (host= user= dbname= sslmode=)
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") || strings.HasPrefix(dsn, "pgx://") ||
		(strings.Contains(dsn, "host=") && (strings.Contains(dsn, "dbname=") || strings.Contains(dsn, "user="))) {
		return "postgres"
	}
	// MySQL: URI scheme or classic DSN like user:pass@tcp(host:port)/db?...
	if strings.HasPrefix(dsn, "mysql://") || strings.Contains(dsn, "@tcp(") || strings.Contains(dsn, ")/") {
		// Heuristic: if it looks like user:pass@tcp(host:port)/db...
		if strings.Contains(dsn, "@tcp(") && strings.Contains(dsn, ")/") {
			return "mysql"
		}
		// mysql:// scheme
		if strings.HasPrefix(dsn, "mysql://") {
			return "mysql"
		}
	}
	// SQLite: file: URIs, :memory:, or *.db path
	if strings.HasPrefix(dsn, "file:") || strings.HasSuffix(dsn, ".db") || strings.HasPrefix(dsn, ":memory:") {
		return "sqlite"
	}
	return ""
}
