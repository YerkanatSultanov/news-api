package repository

import (
	"context"
	"database/sql"
	"errors"
	"news-api/internal/models"
	"news-api/pkg/logger"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password, role, avatar)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.DB.QueryRowContext(
		ctx, query,
		user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.Avatar,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		logger.Log.Error("Error creating user: " + err.Error())
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, first_name, last_name, email, password, role, avatar, created_at FROM users WHERE email=$1`
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.Role, &user.Avatar, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error("Error fetching user by email", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, first_name, last_name, email, password, role, avatar, created_at FROM users WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.Role, &user.Avatar, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Log.Error("Error fetching user by id", err)
		return nil, err
	}
	return user, nil
}
