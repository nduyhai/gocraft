package ports

// All methods must be idempotent.
//
// Parameters:
// - alias:      import alias to use (e.g., "httpgin")
// - importPath: full import path (e.g., "<module>/internal/adapters/inbound/http/gin")
// - optionExpr: expression to include inside fx.Options block (e.g., "httpgin.Module()")
//
// Implementations typically operate relative to a project root path and MUST
// tolerate missing target files without error (no-op behavior).

type AdaptersModuleEditor interface {
	Ensure(alias, importPath, optionExpr string) error
}
