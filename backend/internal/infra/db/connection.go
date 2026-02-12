package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connection holds the database connection pool
type Connection struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection
func New(dsn string) (*Connection, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Connection{Pool: pool}, nil
}

// Close closes the database connection
func (c *Connection) Close() {
	c.Pool.Close()
}
