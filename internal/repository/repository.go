// internal/repository/repository.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL driver
	"github.com/yourusername/azure-go-app/internal/config"
)

type Repository struct {
	db     *sqlx.DB
	config *config.Config
}

func NewRepository(ctx context.Context, cfg *config.Config) (*Repository, error) {
	if cfg.DatabaseURL == "" {
		return &Repository{
			db:     nil, /
			config: cfg,
		}, nil
	}

	db, err := sqlx.ConnectContext(ctx, "sqlserver", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections / 4)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{
		db:     db,
		config: cfg,
	}, nil
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) CheckHealth(ctx context.Context) error {
	if r.db == nil {
		return nil // Skip DB check in local/test mode
	}
	return r.db.PingContext(ctx)
}
func (r *Repository) Close() {
	if r.db != nil {
		_ = r.db.Close()
	}
}