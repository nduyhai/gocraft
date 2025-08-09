package usecase

import (
	"github.com/nduyhai/go-clean-arch-starter/internal/core/ports"
)

// ApplyModules orchestrates applying one or more modules to a given context.
// It delegates to the Registry port, keeping orchestration in the usecase layer.
type ApplyModules struct {
	Registry ports.Registry
}

func (uc ApplyModules) Execute(ctx ports.Ctx, names ...string) error {
	if uc.Registry == nil {
		return nil
	}
	return uc.Registry.Apply(ctx, names...)
}

// ListModules retrieves all registered modules from the Registry port.
type ListModules struct {
	Registry ports.Registry
}

func (uc ListModules) Execute() []ports.Module {
	if uc.Registry == nil {
		return nil
	}
	return uc.Registry.List()
}
