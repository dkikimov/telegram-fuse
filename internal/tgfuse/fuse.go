package tgfuse

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	"telegram-fuse/internal/usecase"
)

type Node struct {
	fs.Inode
	id      int
	storage usecase.Storage
	name    string
}

func NewNode(storage usecase.Storage) *Node {
	return &Node{storage: storage}
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
	if n.GetChild(name) != nil {
		return n.GetChild(name), 0
	}

	return nil, syscall.ENOENT
}

var _ = (fs.NodeOnAdder)((*Node)(nil))

func (n *Node) OnAdd(ctx context.Context) {
	return
}

var _ = (fs.NodeCreater)((*Node)(nil))

func (n *Node) Create(ctx context.Context, name string, _ uint32, _ uint32, out *fuse.EntryOut) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	fileId, err := n.storage.SaveFile(n.id, name, []byte("empty"))
	if err != nil {
		return nil, nil, 0, syscall.EAGAIN
	}

	node = n.NewInode(ctx, n.EmbeddedInode(), defaultAttr)
	fh = NewFile(fileId, n.storage)

	return node, fh, 0, 0
}

var _ = (fs.NodeMkdirer)((*Node)(nil))

func (n *Node) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOSYS
}

var _ = (fs.NodeOpener)((*Node)(nil))

func (n *Node) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, 0, syscall.ENOSYS
}

var _ = (fs.NodeReaddirer)((*Node)(nil))

func (n *Node) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	return nil, syscall.ENOSYS
}

var _ = (fs.NodeReader)((*Node)(nil))

func (n *Node) Read(ctx context.Context, f fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	return nil, syscall.ENOSYS
}

var _ = (fs.NodeRenamer)((*Node)(nil))

func (n *Node) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	return syscall.ENOSYS
}

var _ = (fs.NodeRmdirer)((*Node)(nil))

func (n *Node) Rmdir(ctx context.Context, name string) syscall.Errno {
	return syscall.ENOSYS
}

// var _ = (fs.NodeStatfser())((*Node)(nil))

var _ = (fs.NodeWriter)((*Node)(nil))

func (n *Node) Write(ctx context.Context, f fs.FileHandle, data []byte, off int64) (written uint32, errno syscall.Errno) {
	return 0, syscall.ENOSYS
}

var _ = (fs.NodeGetattrer)((*Node)(nil))

func (n *Node) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	return 0
}
