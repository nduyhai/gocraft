package ports

import "github.com/nduyhai/go-clean-arch-starter/internal/core/entity"

type Renderer interface {
	// Render processes the given template with the context and returns files to write.
	Render(tpl Template, ctx any) ([]entity.File, error)
}
