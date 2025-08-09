package dockerfile

import "embed"

// TemplatesFS embeds the template files for the dockerfile module.
//
//go:embed templates/Dockerfile.tmpl
var TemplatesFS embed.FS
