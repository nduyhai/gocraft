package ports

import "github.com/nduyhai/go-clean-arch-starter/internal/core/entity"

type FSWriter interface {
	WriteAll(root string, files []entity.File) error
}
