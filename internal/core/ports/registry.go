package ports

type Registry interface {
	Register(Module)
	List() []Module
	Get(name string) (Module, bool)
	Apply(ctx Ctx, names ...string) error
}
