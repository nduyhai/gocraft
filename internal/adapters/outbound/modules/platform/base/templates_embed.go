package base

import "embed"

// TemplatesFS embeds the template files for the base module.
// We avoid embedding empty directories (e.g., templates/cmd which only contains
// an "__name__" subdirectory that starts with an underscore and is ignored by default).
// Instead, we explicitly include only the paths that contain embeddable files.
//
// Include top-level template files
// Include subdirectories with actual templates
// Explicitly include the underscored __name__ skeleton
//
//go:embed templates/*.tmpl
//go:embed templates/config/**
//go:embed templates/internal/platform/**
//go:embed templates/cmd/__name__/**
var TemplatesFS embed.FS
