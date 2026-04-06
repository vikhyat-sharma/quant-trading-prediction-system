package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// DBConfig holds database configuration
type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewDB creates a new database connection with optimized settings
func NewDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	config := DBConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// NewDBWithConfig creates a database connection with custom configuration
func NewDBWithConfig(url string, config DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
