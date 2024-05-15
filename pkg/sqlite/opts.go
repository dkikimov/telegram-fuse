package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	// Added for side effect of registering sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	"telegram-fuse/pkg/config"
)

const (
	defaultMaxOpenConns    = 30
	defaultMaxIdleConns    = 15
	defaultMaxConnLifetime = 180
)

// Opts represents options to initializes new sqlite wrapper.
type Opts struct {
	Path            string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime int
}

// New initializes a new sqlite wrapper and verifies that connection is stable.
func New() (*sql.DB, error) {
	// Check db options, set defaults
	opts := &Opts{
		Path:            config.DatabaseCfg.Path,
		MaxOpenConns:    defaultMaxOpenConns,
		MaxIdleConns:    defaultMaxIdleConns,
		MaxConnLifetime: defaultMaxConnLifetime,
	}

	db, err := sql.Open("sqlite3", opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %w", err)
	}

	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(opts.MaxConnLifetime) * time.Second)
	db.SetMaxIdleConns(opts.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
