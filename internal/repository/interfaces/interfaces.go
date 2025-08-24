package interfaces

import (
	"context"
	"news-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
}

type NewsRepository interface {
	Create(news *models.News) error
	Update(news *models.News) error
	Delete(id int) error
	GetByID(id int) (*models.News, error)
	List(params models.NewsListParams) ([]models.News, error)
}
