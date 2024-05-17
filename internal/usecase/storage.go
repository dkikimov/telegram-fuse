package usecase

import "telegram-fuse/internal/entity"

type Storage interface {
	SaveFile(parentId int, name string, data []byte) (entity.FilesystemEntity, error)
	SaveDirectory(parentId int, name string) (entity.FilesystemEntity, error)
	UpdateFile(id int, data []byte) (entity.FilesystemEntity, error)
	UpdateEntity(filesystemEntity entity.FilesystemEntity) (*entity.FilesystemEntity, error)
	ReadFile(id int) ([]byte, error)
	GetDirectoryChildren(id int) ([]entity.FilesystemEntity, error)
	GetEntityById(id int) (entity.FilesystemEntity, error)
}
