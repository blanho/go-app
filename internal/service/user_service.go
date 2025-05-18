// internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/azure-go-app/internal/models"
)

var validate = validator.New()

func (s *Service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	s.tel.TrackEvent("GetUserRequest", map[string]string{
		"id": id,
	}, nil)

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.tel.TrackException(err, map[string]string{
			"operation": "GetUserByID",
			"user_id":   id,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	// Validate input
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, input)
	if err != nil {
		s.tel.TrackException(err, map[string]string{
			"operation": "CreateUser",
		})
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.tel.TrackEvent("UserCreated", map[string]string{
		"user_id": user.ID,
	}, nil)

	return user, nil
}

func (s *Service) ListUsers(ctx context.Context, page, pageSize int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	users, err := s.repo.ListUsers(ctx, pageSize, offset)
	if err != nil {
		s.tel.TrackException(err, map[string]string{
			"operation": "ListUsers",
			"page":      fmt.Sprintf("%d", page),
			"pageSize":  fmt.Sprintf("%d", pageSize),
		})
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}