package models

import "time"

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    int       `json:"author_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NewsListParams struct {
	Limit    int
	Offset   int
	AuthorID *int
	Search   *string
}

func (p *NewsListParams) Normalize() {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

type Actor struct {
	UserID int
	Role   string
}
