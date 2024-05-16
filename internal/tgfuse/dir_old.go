package tgfuse

//
// import (
// 	"github.com/hanwen/go-fuse/v2/fs"
// 	"github.com/hanwen/go-fuse/v2/fuse"
//
// 	"telegram-fuse/internal/entity"
// )
//
// func NewListDirStreamFromEntity(list []entity.FilesystemEntity) fs.DirStream {
// 	fuseDir := make([]fuse.DirEntry, len(list))
// 	for i, el := range list {
// 		fuseDir[i] = fuse.DirEntry{
// 			Mode: defaultAttr.Mode,
// 			Name: el.Name,
// 		}
// 	}
//
// 	return fs.NewListDirStream(fuseDir)
// }
