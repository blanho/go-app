// internal/service/service.go
package service

import (
	"github.com/yourusername/azure-go-app/internal/repository"
	"github.com/yourusername/azure-go-app/internal/telemetry"
)

type Service struct {
	repo *repository.Repository
	tel  *telemetry.Telemetry
}

func NewService(repo *repository.Repository, tel *telemetry.Telemetry) *Service {
	return &Service{
		repo: repo,
		tel:  tel,
	}
}