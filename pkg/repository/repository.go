package repository

import "telegram-fuse/internal/entity"

type Repository interface {
	Close() error
	SaveEntity(filesystemEntity entity.FilesystemEntity) (int, error)
	GetDirectoryChildren(id int) ([]entity.FilesystemEntity, error)
	RenameEntity(filesystemEntity entity.FilesystemEntity) error
	RemoveEntity(filesystemEntity entity.FilesystemEntity) error
	GetEntity(id int) (entity.FilesystemEntity, error)
	UpdateEntity(filesystemEntity entity.FilesystemEntity) error
	DeleteEntity(id int) error
}
