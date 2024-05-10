package usecase

type Storage interface {
	SaveFile(path string, name string, data []byte) error
}
