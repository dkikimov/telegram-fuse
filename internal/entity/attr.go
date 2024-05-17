package entity

import (
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func (f *FilesystemEntity) SetAttr(attr *fuse.Attr) {
	if f.IsDirectory() {
		attr.Mode = fuse.S_IFDIR | 0777
	} else {
		attr.Mode = fuse.S_IFREG | 0777
	}

	attr.Size = uint64(f.Size)
	attr.Ino = uint64(f.Id)
	attr.Atime = uint64(f.CreatedAt.Unix())
	attr.Mtime = uint64(f.UpdatedAt.Unix())
	attr.Ctime = uint64(f.UpdatedAt.Unix())
}

func (f *FilesystemEntity) FromAttr(attr *fuse.SetAttrIn) {
	f.Size = int(attr.Size)
	f.CreatedAt = time.Unix(int64(attr.Atime), int64(attr.Atimensec))
	f.UpdatedAt = time.Unix(int64(attr.Mtime), int64(attr.Mtimensec))
}
