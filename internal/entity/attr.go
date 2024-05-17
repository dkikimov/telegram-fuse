package entity

import (
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func (f *FilesystemEntity) SetAttr(attr *fuse.Attr) {
	attr.Mode = 0777
	attr.Size = uint64(f.Size)
	attr.Ino = uint64(f.Id)
	attr.Atime = uint64(f.CreatedAt.Unix())
	attr.Mtime = uint64(f.UpdatedAt.Unix())
	attr.Ctime = uint64(f.UpdatedAt.Unix())

	attr.Crtime_ = uint64(f.CreatedAt.Unix())
	attr.Crtimensec_ = uint32(f.CreatedAt.Nanosecond())
	attr.Flags_ = 0
}

func (f *FilesystemEntity) FromAttr(attr *fuse.SetAttrIn) {
	f.Size = int(attr.Size)

	f.UpdatedAt = time.Unix(int64(attr.Mtime), int64(attr.Mtimensec))
}
