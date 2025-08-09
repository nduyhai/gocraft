package ginmodule

import "embed"

// TemplatesFS embeds all files under the local templates directory of the gin module.
//
//go:embed templates
var TemplatesFS embed.FS
