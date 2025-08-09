package base

import "embed"

// TemplatesFS embeds all files under the local templates directory of the base module.
// Using the directory name includes all files recursively.
//
//go:embed templates
var TemplatesFS embed.FS
