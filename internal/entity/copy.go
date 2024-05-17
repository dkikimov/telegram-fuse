package entity

func (f *FilesystemEntity) FromEntity(e FilesystemEntity) {
	f.Id = e.Id
	f.ParentId = e.ParentId
	f.Name = e.Name
	f.Size = e.Size
	f.MessageID = e.MessageID
	f.FileID = e.FileID
	f.CreatedAt = e.CreatedAt
	f.UpdatedAt = e.UpdatedAt
}
