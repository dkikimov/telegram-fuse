package entity

import "github.com/hanwen/go-fuse/v2/fuse"

func (f FilesystemEntity) SetAttr(attr *fuse.Attr) bool {
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

	return true
}
