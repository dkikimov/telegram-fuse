package sqlite

type Database struct {
	// ...
}

func NewDatabase() (*Database, error) {
	return &Database{}, nil
}
