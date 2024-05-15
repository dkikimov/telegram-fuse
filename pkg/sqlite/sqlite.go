package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"telegram-fuse/internal/entity"
)

type Database struct {
	*sql.DB
}

func (d *Database) GetEntity(id int) (entity.FilesystemEntity, error) {
	stmt, err := d.DB.Prepare("SELECT id, parent_id, name, file_size, message_id, file_id, created_at, updated_at FROM file_entities WHERE id = ?")
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't prepare get filesystemEntity statement: %w", err)
	}

	var filesystemEntity entity.FilesystemEntity
	err = stmt.QueryRow(id).Scan(&filesystemEntity.Id, &filesystemEntity.ParentId, &filesystemEntity.Name, &filesystemEntity.Size, &filesystemEntity.MessageID, &filesystemEntity.FileID, &filesystemEntity.CreatedAt, &filesystemEntity.UpdatedAt)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't execute get filesystemEntity statement: %w", err)
	}

	return filesystemEntity, nil
}

func (d *Database) RemoveEntity(filesystemEntity entity.FilesystemEntity) error {
	stmt, err := d.DB.Prepare("DELETE FROM file_entities WHERE id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare remove filesystemEntity statement: %w", err)
	}

	_, err = stmt.Exec(filesystemEntity.Id)
	if err != nil {
		return fmt.Errorf("couldn't execute remove filesystemEntity statement: %w", err)
	}

	return nil
}

func (d *Database) SaveEntity(filesystemEntity entity.FilesystemEntity) (int, error) {
	stmt, err := d.DB.Prepare("INSERT INTO file_entities (parent_id, name, file_size, message_id, file_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("couldn't prepare save filesystemEntity statement: %w", err)
	}

	var id int
	err = stmt.QueryRow(filesystemEntity.ParentId, filesystemEntity.Name, filesystemEntity.Size, filesystemEntity.MessageID, filesystemEntity.FileID, filesystemEntity.CreatedAt, filesystemEntity.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("couldn't execute save filesystemEntity statement: %w", err)
	}

	return id, nil
}

func (d *Database) GetDirectoryChildren(filesystemEntity entity.FilesystemEntity) ([]entity.FilesystemEntity, error) {
	if !filesystemEntity.IsDirectory() {
		return nil, fmt.Errorf("filesystemEntity is not a directory")
	}

	stmt, err := d.DB.Prepare("SELECT id, parent_id, name, file_size, message_id, file_id, created_at, updated_at FROM file_entities WHERE parent_id = ?")
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare get directory children statement: %w", err)
	}

	rows, err := stmt.Query(filesystemEntity.Id)
	if err != nil {
		return nil, fmt.Errorf("couldn't execute get directory children statement: %w", err)
	}

	var children []entity.FilesystemEntity
	for rows.Next() {
		var child entity.FilesystemEntity
		err = rows.Scan(&child.Id, &child.ParentId, &child.Name, &child.Size, &child.MessageID, &child.FileID, &child.CreatedAt, &child.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("couldn't scan child filesystemEntity: %w", err)
		}

		children = append(children, child)
	}

	return children, nil
}

func (d *Database) RenameEntity(filesystemEntity entity.FilesystemEntity) error {
	stmt, err := d.DB.Prepare("UPDATE file_entities SET name = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare rename filesystemEntity statement: %w", err)
	}

	_, err = stmt.Exec(filesystemEntity.Name, time.Now(), filesystemEntity.Id)
	if err != nil {
		return fmt.Errorf("couldn't execute rename filesystemEntity statement: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func NewDatabase(DB *sql.DB) *Database {
	return &Database{DB: DB}
}
