package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"news-api/internal/models"
	"news-api/pkg/logger"
	"strings"
)

type NewsRepository struct {
	DB *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{DB: db}
}

func (r *NewsRepository) Create(news *models.News) error {
	query := `
		INSERT INTO news (title, description, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRow(query, news.Title, news.Description, news.AuthorID).
		Scan(&news.ID, &news.CreatedAt, &news.UpdatedAt)
	if err != nil {
		logger.Log.Error("Error creating news", err)
		return err
	}
	return nil
}

func (r *NewsRepository) Update(news *models.News) error {
	query := `
		UPDATE news
		SET title=$1, description=$2, updated_at=NOW()
		WHERE id=$3
		RETURNING updated_at
	`
	err := r.DB.QueryRow(query, news.Title, news.Description, news.ID).Scan(&news.UpdatedAt)
	if err != nil {
		logger.Log.Error("Error updating news", err)
		return err
	}
	return nil
}

func (r *NewsRepository) Delete(id int) error {
	query := `DELETE FROM news WHERE id=$1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		logger.Log.Error("Error deleting news", err)
		return err
	}
	return nil
}

func (r *NewsRepository) GetByID(id int) (*models.News, error) {
	news := &models.News{}
	query := `SELECT id, title, description, author_id, created_at, updated_at FROM news WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(
		&news.ID, &news.Title, &news.Description,
		&news.AuthorID, &news.CreatedAt, &news.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error("Error fetching news by id", err)
		return nil, err
	}
	return news, nil
}

func (r *NewsRepository) List(params models.NewsListParams) ([]models.News, error) {
	query := `SELECT id, title, description, author_id, created_at, updated_at FROM news WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if params.AuthorID != nil {
		query += fmt.Sprintf(" AND author_id=$%d", argPos)
		args = append(args, *params.AuthorID)
		argPos++
	}

	if params.Search != nil && strings.TrimSpace(*params.Search) != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argPos)
		args = append(args, "%"+*params.Search+"%")
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		logger.Log.Error("Error listing news", err)
		return nil, err
	}
	defer rows.Close()

	newsList := []models.News{}
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Description, &n.AuthorID, &n.CreatedAt, &n.UpdatedAt); err != nil {
			logger.Log.Error("Error scanning news row", err)
			continue
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}
