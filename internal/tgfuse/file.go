package tgfuse

import (
	"context"
	"log/slog"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	"telegram-fuse/internal/entity"
	"telegram-fuse/internal/usecase"
)

type File struct {
	id      int
	mu      sync.Mutex
	storage usecase.Storage
	buffer  []byte
}

var _ = (fs.FileHandle)((*File)(nil))
var _ = (fs.FileReleaser)((*File)(nil))
var _ = (fs.FileReader)((*File)(nil))
var _ = (fs.FileAllocater)((*File)(nil))

var _ = (fs.FileFsyncer)((*File)(nil))

func (f *File) Allocate(ctx context.Context, off uint64, size uint64, mode uint32) syscall.Errno {
	return 0
}

func (f *File) Release(ctx context.Context) syscall.Errno {
	return 0
}

func (f *File) Flush(ctx context.Context) (*entity.FilesystemEntity, syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()

	slog.Info("flushing file file", "id", f.id)

	if f.buffer == nil {
		f.buffer = []byte(" ")
	}

	newEntity, err := f.storage.UpdateFile(f.id, f.buffer)
	if err != nil {
		slog.Info("couldn't update file", "error", err, "id", f.id)
		return nil, syscall.EIO
	}

	f.buffer = nil

	return &newEntity, 0
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

func (f *File) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()

	slog.Info("writing file", "id", f.id)

	if f.buffer == nil {
		file, err := f.storage.ReadFile(f.id)
		if err != nil {
			return 0, syscall.EIO
		}

		f.buffer = file
	}

	f.buffer = f.buffer[:off]
	f.buffer = append(f.buffer, data...)

	return uint32(len(data)), 0
}

func NewFile(id int, storage usecase.Storage) *File {
	return &File{id: id, storage: storage}
}
