package ports

type ConfigEditor interface {
	// EnsureDefaultsFor updates config/config.yml to include default properties for the given module.
	// It must be idempotent and tolerant to missing files; when the config file is absent, it should create it.
	EnsureDefaultsFor(module string) error
}
