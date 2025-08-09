package chimodule

import "embed"

// TemplatesFS embeds all files under the local templates directory of the chi module.
//
//go:embed templates
var TemplatesFS embed.FS
