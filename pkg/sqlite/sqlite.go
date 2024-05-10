package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	// Added for side effect of registering sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
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
func New(opts *Opts) (*sql.DB, error) {
	// Check db options, set defaults
	if opts.MaxOpenConns == 0 {
		opts.MaxOpenConns = defaultMaxOpenConns
	}
	if opts.MaxIdleConns == 0 {
		opts.MaxIdleConns = defaultMaxIdleConns
	}
	if opts.MaxConnLifetime == 0 {
		opts.MaxConnLifetime = defaultMaxConnLifetime
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
