package entity

import (
	"log/slog"
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
)

func (f *FilesystemEntity) SetEntryOut(entry *fuse.EntryOut) {
	f.SetAttr(&entry.Attr)

	entry.Ino = uint64(f.Id)
}

func (f *FilesystemEntity) SetAttr(attr *fuse.Attr) {
	if f.IsDirectory() {
		attr.Mode = fuse.S_IFDIR | 0755
		attr.Nlink = 2
	} else {
		attr.Mode = fuse.S_IFREG | 0777
		attr.Nlink = 1
	}

	slog.Info("SetAttr", "size", f.Size, "id", f.Id, "created_at", f.CreatedAt, "updated_at", f.UpdatedAt)

	attr.Size = uint64(f.Size)
	attr.Ino = uint64(f.Id)

	attr.Atime = uint64(f.CreatedAt.Unix())
	attr.Mtime = uint64(f.UpdatedAt.Unix())
	attr.Ctime = uint64(f.UpdatedAt.Unix())
	attr.Atimensec = uint32(f.CreatedAt.Nanosecond())
	attr.Mtimensec = uint32(f.UpdatedAt.Nanosecond())
	attr.Ctimensec = uint32(f.UpdatedAt.Nanosecond())
	attr.Blocks = 1

	attr.Crtime_ = uint64(f.CreatedAt.Unix())
	attr.Crtimensec_ = uint32(f.CreatedAt.Nanosecond())
	attr.Flags_ = 0
}

func (f *FilesystemEntity) FromAttr(attr *fuse.SetAttrIn) {
	f.Size = int(attr.Size)

	f.UpdatedAt = time.Unix(int64(attr.Mtime), int64(attr.Mtimensec))
}
