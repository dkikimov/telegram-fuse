package fuse

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type TgNode struct {
	fs.Inode
}

var _ = (fs.InodeEmbedder)((*TgNode)(nil))

var _ = (fs.NodeLookuper)((*TgNode)(nil))

func (n *TgNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	ops := TgNode{}
	out.Mode = 0755
	out.Size = 42
	return n.NewInode(ctx, &ops, fs.StableAttr{Mode: syscall.S_IFREG}), 0
}
