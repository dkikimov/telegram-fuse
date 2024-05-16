package tgfuse

import (
	"context"
	"log/slog"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	"telegram-fuse/internal/usecase"
)

type Node struct {
	fs.Inode
	Id       int
	RootData *Root
	Name     string
	storage  usecase.Storage
}

var defaultAttr = fs.StableAttr{
	Mode: 0777,
}

var _ = (fs.InodeEmbedder)((*Node)(nil))

var _ = (fs.NodeAccesser)((*Node)(nil))

func (n *Node) Access(ctx context.Context, mask uint32) syscall.Errno {
	return syscall.F_OK
}

var _ = (fs.NodeLookuper)((*Node)(nil))

func (n *Node) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	if n.IsDir() == false {
		if name == n.Name {
			return n.NewInode(ctx, n.EmbeddedInode(), defaultAttr), 0
		} else {
			return nil, syscall.ENOENT
		}
	}

	// Check if the file exists in the directory
	fileId, err := n.storage.GetDirectoryChildren(n.Id)
	if err != nil {
		slog.Info("failed to get directory children", "error", err)
		return nil, syscall.EAGAIN
	}

	for _, file := range fileId {
		if file.Name == name {
			node := n.RootData.newNode(file)
			ch := n.NewInode(ctx, node, defaultAttr)

			return ch, 0
		}
	}

	return nil, syscall.ENOENT
}

var _ = (fs.NodeOnAdder)((*Node)(nil))

func (n *Node) OnAdd(ctx context.Context) {
	return
}

var _ = (fs.NodeCreater)((*Node)(nil))
var _ = fs.NodeCreater(&Node{})

func (n *Node) Create(ctx context.Context, name string, _ uint32, _ uint32, out *fuse.EntryOut) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {
	filesystemEntity, err := n.storage.SaveFile(n.Id, name, []byte("empty"))
	if err != nil {
		return nil, nil, 0, syscall.EAGAIN
	}

	node := n.RootData.newNode(filesystemEntity)
	ch := n.NewInode(ctx, node, defaultAttr)

	fh := NewFile(filesystemEntity.Id, n.storage)

	return ch, fh, 0, 0
}

var _ = (fs.NodeMkdirer)((*Node)(nil))

func (n *Node) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOSYS
}

var _ = (fs.NodeOpener)((*Node)(nil))

func (n *Node) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	fh = NewFile(2, n.storage)
	return fh, 0, 0
}

var _ = (fs.NodeReaddirer)((*Node)(nil))

func (n *Node) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	if n.IsDir() == false {
		return nil, syscall.ENOTDIR
	}

	ent, err := n.storage.GetDirectoryChildren(n.Id)
	if err != nil {
		slog.Info("failed to get directory children", "error", err)
		return nil, syscall.EAGAIN
	}

	return NewListDirStreamFromEntity(ent), 0
}

var _ = (fs.NodeRenamer)((*Node)(nil))

func (n *Node) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	return syscall.ENOSYS
}

var _ = (fs.NodeRmdirer)((*Node)(nil))

func (n *Node) Rmdir(ctx context.Context, name string) syscall.Errno {
	return syscall.ENOSYS
}

var _ = (fs.NodeWriter)((*Node)(nil))

func (n *Node) Write(ctx context.Context, f fs.FileHandle, data []byte, off int64) (written uint32, errno syscall.Errno) {
	return 0, syscall.ENOSYS
}

var _ = (fs.NodeGetattrer)((*Node)(nil))

func (n *Node) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	filesystemEntity, err := n.storage.GetEntityById(n.Id)
	if err != nil {
		slog.Info("failed to get entity by id", "error", err)
		return syscall.EAGAIN
	}

	filesystemEntity.SetAttr(&out.Attr)

	return 0
}

var _ = (fs.NodeStatfser)((*Node)(nil))

func (n *Node) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	out.Blocks = 1
	out.Bfree = 1
	out.Bavail = 1

	return 0
}

var _ = (fs.NodeSetattrer)((*Node)(nil))

func (n *Node) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	out.Mode = in.Mode
	out.Size = in.Size
	return 0
}
