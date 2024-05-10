package entity

import "time"

// FilesystemEntity represents a file or directory in the system.
type FilesystemEntity struct {
	// ID is the unique identifier of the entity.
	Id uint64

	// ParentID is the unique identifier of the parent directory.
	ParentId uint64

	// Name is the name of the entity. Only for file.
	Name string

	// Size is the size of the entity in bytes. Only for file.
	Size int

	// MessageID is the unique identifier of the message in the chat. Only for file.
	MessageID int

	// FileID is the unique identifier of the file in the telegram chat.
	// It is used to download the file from the chat.
	FileID string

	// CreatedAt is the time when the file was created.
	CreatedAt time.Time

	// UpdatedAt is the time when the file was last updated.
	UpdatedAt time.Time
}

func NewFile(id, parentId uint64, name string, size int, messageID int, fileID string, createdAt, updatedAt time.Time) FilesystemEntity {
	return FilesystemEntity{
		Id:        id,
		ParentId:  parentId,
		Name:      name,
		Size:      size,
		MessageID: messageID,
		FileID:    fileID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewDirectory(id, parentId uint64, name string, createdAt, updatedAt time.Time) FilesystemEntity {
	return FilesystemEntity{
		Id:        id,
		ParentId:  parentId,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
