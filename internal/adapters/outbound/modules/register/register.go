package register

import (
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/modules/platform/base"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/ports"
)

// Builtins registers all built-in modules into the provided registry.
// Usage:
//
//	r := simple.New()
//	register.Builtins(r)
//	// r now contains platform:base and others when added in future
func Builtins(r ports.Registry) {
	// Platform base module (Fx + Viper, logger, DI root, basic structure)
	r.Register(base.New())
	// TODO: add other built-in modules here as they are implemented, e.g. http:gin, grpc:server, db:postgres, etc.
}
