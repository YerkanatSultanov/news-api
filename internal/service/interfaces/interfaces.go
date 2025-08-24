package interfaces

import (
	"context"
	"news-api/internal/dto/auth"
	"news-api/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, input auth.RegisterUserInput) error
	Login(ctx context.Context, input auth.LoginUserInput) (string, string, error)
}

type NewsService interface {
	CreateNews(ctx context.Context, actor models.Actor, n *models.News) error
	UpdateNews(ctx context.Context, actor models.Actor, n *models.News) error
	DeleteNews(ctx context.Context, actor models.Actor, id int) error
	GetByIDNews(ctx context.Context, id int) (*models.News, error)
	ListNews(ctx context.Context, p models.NewsListParams) ([]models.News, error)
}
