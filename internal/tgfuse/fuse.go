package tgfuse

import (
	"context"
	"log/slog"
	"math"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	"telegram-fuse/internal/entity"
	"telegram-fuse/internal/usecase"
)

type Node struct {
	fs.Inode
	entity.FilesystemEntity
	RootData *Root
	storage  usecase.Storage
}

var _ = (fs.InodeEmbedder)((*Node)(nil))
var _ = (fs.NodeAccesser)((*Node)(nil))

func (n *Node) Access(ctx context.Context, mask uint32) syscall.Errno {
	return syscall.F_OK
}

var _ = (fs.NodeLookuper)((*Node)(nil))

func (n *Node) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	// Check if looking for itself
	if n.IsDirectory() == false {
		if n.Name == name {
			n.SetEntryOut(out)
			return n.EmbeddedInode(), 0
		} else {
			return nil, syscall.ENOENT
		}
	}

	// Check if the file exists in the directory
	filesystemEntities, err := n.storage.GetDirectoryChildren(n.Id)
	if err != nil {
		slog.Info("failed to get directory children", "error", err)
		return nil, syscall.EAGAIN
	}

	for _, file := range filesystemEntities {
		if file.Name == name {
			node := n.RootData.newNode(file)
			ch := n.NewInode(ctx, node, node.GetStableAttr())

			node.SetEntryOut(out)
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

func (n *Node) Create(ctx context.Context, name string, _ uint32, _ uint32, out *fuse.EntryOut) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {
	filesystemEntity, err := n.storage.SaveFile(n.Id, name, []byte(" "))
	if err != nil {
		return nil, nil, 0, syscall.EAGAIN
	}

	node := n.RootData.newNode(filesystemEntity)
	ch := n.NewInode(ctx, node, node.GetStableAttr())

	fh := NewFile(filesystemEntity.Id, n.storage)

	node.SetEntryOut(out)

	slog.Info("created file", "name", name, "id", filesystemEntity.Id, "messageId", filesystemEntity.MessageID)
	return ch, fh, 0, 0
}

var _ = (fs.NodeMkdirer)((*Node)(nil))

func (n *Node) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	directoryEntity, err := n.storage.SaveDirectory(n.Id, name)
	if err != nil {
		slog.Info("failed to save directory", "error", err)
		return nil, syscall.EAGAIN
	}

	node := n.RootData.newNode(directoryEntity)
	ch := n.NewInode(ctx, node, node.GetStableAttr())

	node.SetAttr(&out.Attr)

	return ch, 0
}

var _ = (fs.NodeOpener)((*Node)(nil))

func (n *Node) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	fh = NewFile(n.Id, n.storage)
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
	parentNode := newParent.(*Node)

	children, err := n.storage.GetDirectoryChildren(parentNode.Id)
	if err != nil {
		slog.Info("failed to get directory children", "error", err)
		return syscall.EAGAIN
	}

	for _, child := range children {
		if child.Name == newName {
			return syscall.EEXIST
		}
	}

	nodeToRename := n.GetChild(name).Operations().(*Node)

	nodeToRename.Name = newName
	nodeToRename.ParentId = parentNode.Id

	newEntity, err := n.storage.UpdateEntity(nodeToRename.FilesystemEntity)
	if err != nil {
		slog.Info("failed to update entity", "error", err)
		return syscall.EAGAIN
	}

	nodeToRename.FromEntity(*newEntity)

	slog.Info("renamed entity", "oldName", name, "newName", newName)
	return 0
}

var _ = (fs.NodeUnlinker)((*Node)(nil))

func (n *Node) Unlink(ctx context.Context, name string) syscall.Errno {
	nodeToDelete := n.GetChild(name).Operations().(*Node)

	err := n.storage.DeleteEntity(nodeToDelete.Id)
	if err != nil {
		slog.Info("failed to delete entity", "error", err)
		return syscall.EAGAIN
	}

	slog.Info("deleted entity", "name", name)
	return 0
}

var _ = (fs.NodeRmdirer)((*Node)(nil))

func (n *Node) Rmdir(ctx context.Context, name string) syscall.Errno {
	nodeToDelete := n.GetChild(name).Operations().(*Node)

	err := n.storage.DeleteEntity(nodeToDelete.Id)
	if err != nil {
		slog.Info("failed to delete entity", "error", err)
		return syscall.EAGAIN
	}

	slog.Info("deleted entity", "name", name)
	return 0
}

var _ = (fs.NodeGetattrer)((*Node)(nil))

func (n *Node) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	n.SetAttr(&out.Attr)

	return 0
}

var _ = (fs.NodeSetattrer)((*Node)(nil))

func (n *Node) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	n.SetAttr(&out.Attr)

	return 0
}

var _ = (fs.NodeWriter)((*Node)(nil))

func (n *Node) Write(ctx context.Context, f fs.FileHandle, data []byte, off int64) (written uint32, errno syscall.Errno) {
	file, ok := f.(*File)
	if !ok {
		return 0, syscall.EINVAL
	}

	written, e := file.Write(ctx, data, off)
	if e != 0 {
		return 0, fs.ToErrno(e)
	}

	slog.Info("wrote to file", "id", n.Id, "size", n.Size, "messageId", n.MessageID)

	return written, 0
}

var _ = (fs.NodeFlusher)((*Node)(nil))

func (n *Node) Flush(ctx context.Context, f fs.FileHandle) syscall.Errno {
	file, ok := f.(*File)
	if !ok {
		return syscall.EINVAL
	}

	newEntity, e := file.Flush(ctx)
	if e != 0 {
		slog.Error("failed to flush file", "error", e)
		return fs.ToErrno(e)
	}

	n.FromEntity(*newEntity)

	slog.Info("flushed file", "id", n.Id, "size", n.Size, "messageId", n.MessageID)

	return 0
}

var _ = (fs.NodeStatfser)((*Node)(nil))

var basicStatfs = syscall.Statfs_t{
	Bsize:  4096,
	Iosize: math.MaxInt32,
	Blocks: math.MaxUint64,
	Bfree:  math.MaxUint64,
	Bavail: math.MaxUint64,
	Files:  math.MaxUint64,
	Ffree:  math.MaxUint64,
}

func (n *Node) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	out.FromStatfsT(&basicStatfs)

	return 0
}
