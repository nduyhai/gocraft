package register

import (
	dockerfilemodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/feature/dockerfile"
	gitignoremodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/feature/gitignore"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/modules/feature/makefile"
	grpcservermodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/grpc/server"
	chimodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/http/chi"
	ginmodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/http/gin"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/modules/platform/base"
	"github.com/nduyhai/gocraft/internal/core/ports"
)

// Builtins registers all built-in modules into the provided registry.
// Usage:
//
//	r := embed_registry.New()
//	register.Builtins(r)
//	// r now contains platform:base and others when added in future
func Builtins(r ports.Registry) {
	// Platform base module (Fx + Viper, logger, DI root, basic structure)
	r.Register(base.New())
	// HTTP server modules
	r.Register(ginmodule.New())
	r.Register(chimodule.New())
	// gRPC server module
	r.Register(grpcservermodule.New())
	// Feature modules
	r.Register(gitignoremodule.New())
	r.Register(makefile.New())
	r.Register(dockerfilemodule.New())
}
