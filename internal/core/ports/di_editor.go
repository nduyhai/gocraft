package ports

type DependencyInjectionEditor interface {
	Ensure(alias, importPath, optionExpr string) error
}
