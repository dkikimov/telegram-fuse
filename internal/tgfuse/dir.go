package tgfuse

import (
	"context"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"telegram-fuse/internal/usecase"
)

type Dir struct {
	Id      int
	storage usecase.Storage
}

var _ fs.Node = (*Dir)(nil)

func (d *Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = 0755

	return nil
}

var _ fs.NodeRequestLookuper = (*Dir)(nil)

func (d *Dir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	children, err := d.storage.GetDirectoryChildren(d.Id)
	if err != nil {
		return nil, syscall.EIO
	}

	for _, child := range children {
		if child.Name == req.Name {
			if child.IsDirectory() {
				return &Dir{
					Id:      child.Id,
					storage: d.storage,
				}, nil
			}

			return &File{
				Id:      child.Id,
				storage: d.storage,
			}, nil
		}
	}

	return nil, syscall.ENOENT
}

var _ fs.HandleReadDirAller = (*Dir)(nil)

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	children, err := d.storage.GetDirectoryChildren(d.Id)
	if err != nil {
		return nil, syscall.EDEVERR
	}

	var dirents []fuse.Dirent
	for _, child := range children {
		var dirent fuse.Dirent
		if child.IsDirectory() {
			dirent.Type = fuse.DT_Dir
		} else {
			dirent.Type = fuse.DT_File
		}
		dirent.Name = child.Name

		dirents = append(dirents, dirent)
	}

	return dirents, nil
}
