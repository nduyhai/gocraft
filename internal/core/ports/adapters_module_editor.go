package ports

// AdaptersModuleEditor provides operations to update the generated project's
// internal/adapters/module.go file in an idempotent manner.
// It ensures an import alias/path exists and that an fx option expression
// (e.g., httpgin.Module()) is included in the returned fx.Options(...).
//
// Implementations should be resilient to different import block styles
// (grouped, single-line, or absent) and insert the option line either
// after a known needle or before the closing parenthesis of the fx.Options block.
//
// All methods should be idempotent.
//
// alias:      import alias to use (e.g., "httpgin")
// importPath: full import path (e.g., "<module>/internal/adapters/inbound/http/gin")
// optionExpr: expression to include inside fx.Options block (e.g., "httpgin.Module()")
//
// Implementations typically operate relative to a project root path.
//
// Note: This editor intentionally targets the base template's adapters module
// structure. If that template changes fundamentally, the implementation may
// need adjustments.

type AdaptersModuleEditor interface {
	Ensure(alias, importPath, optionExpr string) error
}
