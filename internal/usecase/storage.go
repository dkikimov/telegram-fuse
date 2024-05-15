package usecase

type Storage interface {
	SaveFile(id uint64, parentId uint64, name string, data []byte) error
	ReadFile(id uint64) ([]byte, error)
}
