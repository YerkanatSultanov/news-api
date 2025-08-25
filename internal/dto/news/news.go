package news

type News struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description" binding:"required"`
}

type UpdateNewsRequest struct {
	ID          int     `json:"id" binding:"required"`
	Title       *string `json:"title,omitempty" binding:"max=255"`
	Description *string `json:"description,omitempty"`
}
