package tgfuse

//
// import (
// 	"context"
// 	"sync"
// 	"syscall"
//
// 	"github.com/hanwen/go-fuse/v2/fs"
// 	"github.com/hanwen/go-fuse/v2/fuse"
//
// 	"telegram-fuse/internal/usecase"
// )
//
// type File struct {
// 	id      int
// 	mu      sync.Mutex
// 	storage usecase.Storage
// }
//
// func NewFile(id int, storage usecase.Storage) *File {
// 	return &File{id: id, storage: storage}
// }
//
// var _ = (fs.FileHandle)((*File)(nil))
//
// // var _ = (fs.FileReleaser)((*File)(nil))
// // var _ = (fs.FileGetattrer)((*File)(nil))
// var _ = (fs.FileReader)((*File)(nil))
//
// // var _ = (fs.FileWriter)((*File)(nil))
// // var _ = (fs.FileGetlker)((*File)(nil))
// // var _ = (fs.FileSetlker)((*File)(nil))
// // var _ = (fs.FileSetlkwer)((*File)(nil))
// // var _ = (fs.FileLseeker)((*File)(nil))
// // var _ = (fs.FileFlusher)((*File)(nil))
// // var _ = (fs.FileFsyncer)((*File)(nil))
// // var _ = (fs.FileSetattrer)((*File)(nil))
// // var _ = (fs.FileAllocater)((*File)(nil))
//
// func (f *File) Read(ctx context.Context, buf []byte, off int64) (res fuse.ReadResult, errno syscall.Errno) {
// 	f.mu.Lock()
// 	defer f.mu.Unlock()
//
// 	file, err := f.storage.ReadFile(f.id)
// 	if err != nil {
// 		return nil, syscall.EIO
// 	}
//
// 	res = fuse.ReadResultData(file)
//
// 	return res, syscall.F_OK
// }
