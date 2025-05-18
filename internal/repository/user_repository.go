// internal/repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/azure-go-app/internal/models"
)

func (r *Repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if r.db == nil {
		return nil, errors.New("database not connected")
	}

	user := &models.User{}
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = @id`
	
	err := r.db.GetContext(ctx, user, query, sql.Named("id", id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil 
		}
		return nil, err
	}
	
	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	if r.db == nil {
		return nil, errors.New("database not connected")
	}

	now := time.Now()
	user := &models.User{
		ID:        uuid.New().String(),
		Username:  input.Username,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	query := `
		INSERT INTO users (id, username, email, created_at, updated_at)
		VALUES (@id, @username, @email, @created_at, @updated_at)
	`
	
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *Repository) ListUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	if r.db == nil {
		return nil, errors.New("database not connected")
	}

	var users []models.User
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		OFFSET @offset ROWS
		FETCH NEXT @limit ROWS ONLY
	`
	
	err := r.db.SelectContext(ctx, &users, query, sql.Named("offset", offset), sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	
	return users, nil
}