package ports

type GoModEditor interface {
	Add(module, version string) error
	Replace(oldPath, newPath string) error
	Tidy() error
}
