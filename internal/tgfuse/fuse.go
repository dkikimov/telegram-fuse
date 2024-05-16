package tgfuse

import (
	"bazil.org/fuse/fs"

	"telegram-fuse/internal/usecase"
)

type Filesystem struct {
	Id      int
	Storage usecase.Storage
}

var _ fs.FS = (*Filesystem)(nil)

func (f *Filesystem) Root() (fs.Node, error) {
	n := &Dir{
		Id:      0,
		storage: f.Storage,
	}

	return n, nil
}
