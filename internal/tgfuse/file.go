package tgfuse

import (
	"context"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"telegram-fuse/internal/usecase"
)

type File struct {
	Id      int
	storage usecase.Storage
}

var _ fs.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = 0755

	return nil
}
