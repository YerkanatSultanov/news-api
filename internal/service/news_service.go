package service

import (
	"context"
	"errors"
	errors2 "news-api/internal/dto/errors"
	"strings"
	"unicode/utf8"

	"news-api/internal/models"
	"news-api/internal/repository/interfaces"
	"news-api/pkg/logger"
)

type NewsService struct {
	repo interfaces.NewsRepository
}

func NewNewsService(repo interfaces.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

const (
	maxTitleLen = 255
	adminRole   = "admin"
	editorRole  = "editor"
)

func (s *NewsService) CreateNews(ctx context.Context, actor models.Actor, n *models.News) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if actor.Role != adminRole && actor.Role != editorRole {
		logger.Log.Warn("Create news forbidden: role mismatch", "role", actor.Role)
		return errors2.ErrForbidden
	}

	if err := validateNewsPayload(n); err != nil {
		logger.Log.Warn("Create news validation failed", "error", err)
		return err
	}

	n.AuthorID = actor.UserID

	if err := s.repo.Create(n); err != nil {
		logger.Log.Error("Create news failed", "error", err, "author_id", n.AuthorID)
		return err
	}

	logger.Log.Info("News created", "news_id", n.ID, "author_id", n.AuthorID)
	return nil
}

func (s *NewsService) UpdateNews(ctx context.Context, actor models.Actor, n *models.News) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if actor.Role != adminRole && actor.Role != editorRole {
		logger.Log.Warn("Update news forbidden: role mismatch", "role", actor.Role)
		return errors2.ErrForbidden
	}
	if n.ID <= 0 {
		return errors.Join(errors2.ErrValidation, errors.New("id is required"))
	}

	if err := validateNewsPayload(n); err != nil {
		logger.Log.Warn("Update news validation failed", "error", err, "news_id", n.ID)
		return err
	}

	existing, err := s.repo.GetByID(n.ID)
	if err != nil {
		logger.Log.Error("GetByID before update failed", "error", err, "news_id", n.ID)
		return err
	}
	if existing == nil {
		logger.Log.Warn("News not found for update", "news_id", n.ID)
		return errors2.ErrNotFound
	}

	if err := s.repo.Update(n); err != nil {
		logger.Log.Error("Update news failed", "error", err, "news_id", n.ID)
		return err
	}

	logger.Log.Info("News updated", "news_id", n.ID)
	return nil
}

func (s *NewsService) DeleteNews(ctx context.Context, actor models.Actor, id int) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if actor.Role != adminRole {
		logger.Log.Warn("Delete news forbidden: non-admin", "role", actor.Role)
		return errors2.ErrForbidden
	}
	if id <= 0 {
		return errors.Join(errors2.ErrValidation, errors.New("id is required"))
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		logger.Log.Error("GetByID before delete failed", "error", err, "news_id", id)
		return err
	}
	if existing == nil {
		logger.Log.Warn("News not found for delete", "news_id", id)
		return errors2.ErrNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		logger.Log.Error("Delete news failed", "error", err, "news_id", id)
		return err
	}

	logger.Log.Info("News deleted", "news_id", id)
	return nil
}

func (s *NewsService) GetByIDNews(ctx context.Context, id int) (*models.News, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if id <= 0 {
		return nil, errors.Join(errors2.ErrValidation, errors.New("id is required"))
	}

	n, err := s.repo.GetByID(id)
	if err != nil {
		logger.Log.Error("GetByID failed", "error", err, "news_id", id)
		return nil, err
	}
	if n == nil {
		return nil, errors2.ErrNotFound
	}
	return n, nil
}

func (s *NewsService) ListNews(ctx context.Context, p models.NewsListParams) ([]models.News, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	p.Normalize()
	list, err := s.repo.List(p)
	if err != nil {
		logger.Log.Error("List news failed", "error", err, "limit", p.Limit, "offset", p.Offset)
		return nil, err
	}
	logger.Log.Info("List news ok", "count", len(list), "limit", p.Limit, "offset", p.Offset)
	return list, nil
}

func validateNewsPayload(n *models.News) error {
	title := strings.TrimSpace(n.Title)
	desc := strings.TrimSpace(n.Description)

	if title == "" || desc == "" {
		return errors.Join(errors2.ErrValidation, errors.New("title and description are required"))
	}
	if utf8.RuneCountInString(title) > maxTitleLen {
		return errors.Join(errors2.ErrValidation, errors.New("title exceeds 255 chars"))
	}
	return nil
}
