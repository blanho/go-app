// internal/repository/health_repository.go
package repository

import (
	"context"
)

type HealthStatus struct {
	Database bool `json:"database"`
	Redis    bool `json:"redis"`
}

func (r *Repository) HealthStatus(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Database: true,
		Redis:    true,
	}

	if r.db != nil {
		err := r.db.PingContext(ctx)
		status.Database = err == nil
	}

	return status
}