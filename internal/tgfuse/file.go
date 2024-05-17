package tgfuse

import (
	"context"
	"log/slog"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	"telegram-fuse/internal/usecase"
)

type File struct {
	id      int
	mu      sync.Mutex
	storage usecase.Storage
}

var _ = (fs.FileHandle)((*File)(nil))
var _ = (fs.FileReleaser)((*File)(nil))
var _ = (fs.FileReader)((*File)(nil))
var _ = (fs.FileFlusher)((*File)(nil))
var _ = (fs.FileAllocater)((*File)(nil))

var _ = (fs.FileFsyncer)((*File)(nil))

func (f *File) Allocate(ctx context.Context, off uint64, size uint64, mode uint32) syscall.Errno {
	return 0
}

func (f *File) Release(ctx context.Context) syscall.Errno {
	return 0
}

func (f *File) Flush(ctx context.Context) syscall.Errno {
	return 0
}

func (f *File) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	return 0
}

func (f *File) Read(ctx context.Context, buf []byte, off int64) (res fuse.ReadResult, errno syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()

	slog.Info("reading file", "id", f.id)

	file, err := f.storage.ReadFile(f.id)
	if err != nil {
		return nil, syscall.EIO
	}

	res = fuse.ReadResultData(file)

	return res, syscall.F_OK
}

func (f *File) Write(ctx context.Context, data []byte, off int64) (total int, written uint32, errno syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()

	slog.Info("writing file", "id", f.id)

	content, err := f.storage.ReadFile(f.id)
	if err != nil {
		return 0, 0, syscall.EIO
	}

	content = content[:off]
	content = append(content, data...)

	_, err = f.storage.UpdateFile(f.id, content)
	if err != nil {
		return 0, 0, syscall.EIO
	}

	return len(content), uint32(len(data)), 0
}

func NewFile(id int, storage usecase.Storage) *File {
	return &File{id: id, storage: storage}
}
