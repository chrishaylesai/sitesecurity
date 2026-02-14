package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/chrishaylesai/sitesecurity/api/internal/config"
)

// NewDB creates a new database connection.
func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
