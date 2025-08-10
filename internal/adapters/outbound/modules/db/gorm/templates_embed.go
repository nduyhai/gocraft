package gormmodule

import "embed"

// TemplatesFS embeds all files under the local templates directory of the gorm module.
//
//go:embed templates
var TemplatesFS embed.FS
