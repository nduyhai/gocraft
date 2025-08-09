package ports

import "github.com/nduyhai/gocraft/internal/core/entity"

type FSWriter interface {
	WriteAll(root string, files []entity.File) error
}
