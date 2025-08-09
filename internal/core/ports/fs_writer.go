package ports

import "github.com/you/cleanctl/internal/core/entity"

type FSWriter interface {
	WriteAll(root string, files []entity.File) error
}
