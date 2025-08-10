package cli

import (
	configfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/config/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/context/contextimpl"
	amfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/di/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/fs/oswriter"
	gomodfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/gomod/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/rendering/texttmpl"
	"github.com/nduyhai/gocraft/internal/core/ports"
)

// newCtx builds a ports.Ctx with common outbound collaborators.
func newCtx(root, name, module string) ports.Ctx {
	renderer := texttmpl.New()
	writer := oswriter.New()
	gomod := gomodfileeditor.New(root)
	adaptersEditor := amfileeditor.New(root)
	cfgEditor := configfileeditor.New(root)

	return contextimpl.New(root, writer, renderer, gomod, adaptersEditor, cfgEditor, map[string]any{
		"Name":   name,
		"Module": module,
	})
}
