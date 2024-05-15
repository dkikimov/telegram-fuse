package usecase

type Storage interface {
	SaveFile(parentId int, name string, data []byte) (int, error)
	ReadFile(id int) ([]byte, error)
}
