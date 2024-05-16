package usecase

import "telegram-fuse/internal/entity"

type Storage interface {
	SaveFile(parentId int, name string, data []byte) (entity.FilesystemEntity, error)
	ReadFile(id int) ([]byte, error)
	GetDirectoryChildren(id int) ([]entity.FilesystemEntity, error)
	GetEntityById(id int) (entity.FilesystemEntity, error)
}
