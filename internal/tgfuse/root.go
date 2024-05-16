package tgfuse

import (
	"time"

	"github.com/hanwen/go-fuse/v2/fs"

	"telegram-fuse/internal/entity"
	"telegram-fuse/internal/usecase"
)

type Root struct {
	storage usecase.Storage
}

func (r *Root) newNode(filesystemEntity entity.FilesystemEntity) fs.InodeEmbedder {
	return &Node{
		RootData: r,
		storage:  r.storage,
		Id:       filesystemEntity.Id,
		Name:     filesystemEntity.Name,
	}
}

func NewRoot(storage usecase.Storage) fs.InodeEmbedder {
	root := &Root{
		storage: storage,
	}

	return root.newNode(entity.NewDirectory(
		0,
		-1,
		"root",
		time.Now(),
		time.Now(),
	))
}
