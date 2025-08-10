package grpcservermodule

import "embed"

// TemplatesFS embeds all files under the local templates directory of the grpc server module.
//
//go:embed templates
var TemplatesFS embed.FS
