package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"news-api/internal/middleware"
	"strconv"

	"github.com/gorilla/mux"
	errors2 "news-api/internal/dto/errors"
	"news-api/internal/dto/news"
	"news-api/internal/models"
	"news-api/internal/service/interfaces"
	"news-api/pkg/logger"
	"news-api/utils"
)

type NewsHandler struct {
	newsService interfaces.NewsService
}

func NewNewsHandler(newsService interfaces.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

// ListNews godoc
// @Summary      Get news list
// @Description  Returns list of news with pagination, filtering and search
// @Tags         news
// @Produce      json
// @Param        limit     query   int     false  "Limit (default 10)"
// @Param        offset    query   int     false  "Offset (default 0)"
// @Param        author_id query   int     false  "Filter by author id"
// @Param        search    query   string  false  "Search by title"
// @Success      200  {array}   models.News
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/news [get]
func (h *NewsHandler) ListNews(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	params := models.NewsListParams{
		Limit:  limit,
		Offset: offset,
	}

	if authorStr := r.URL.Query().Get("author_id"); authorStr != "" {
		if authorID, err := strconv.Atoi(authorStr); err == nil {
			params.AuthorID = &authorID
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		params.Search = &search
	}

	news, err := h.newsService.ListNews(r.Context(), params)
	if err != nil {
		logger.Log.Error("list news failed", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to list news")
		return
	}

	utils.WriteJSON(w, http.StatusOK, news)
}

// GetNewsByID godoc
// @Summary      Get news by ID
// @Tags         news
// @Produce      json
// @Param        id   path   int  true  "News ID"
// @Success      200  {object}  models.News
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/news/{id} [get]
func (h *NewsHandler) GetNewsByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	n, err := h.newsService.GetByIDNews(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors2.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, "news not found")
			return
		}
		logger.Log.Error("get news failed", "error", err, "id", id)
		utils.WriteError(w, http.StatusInternalServerError, "failed to get news")
		return
	}

	utils.WriteJSON(w, http.StatusOK, n)
}

// CreateNews godoc
// @Summary      Create news
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        input  body   news.News  true  "News input"
// @Success      201  {object}  news.News
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Security     BearerAuth
// @Router       /api/news [post]
func (h *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	var n news.News
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	actor, ok := getActor(r)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	newsTemp := models.News{
		Title:       n.Title,
		Description: n.Description,
	}

	if err := h.newsService.CreateNews(r.Context(), actor, &newsTemp); err != nil {
		switch {
		case errors.Is(err, errors2.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, "forbidden")
		case errors.Is(err, errors2.ErrValidation):
			utils.WriteError(w, http.StatusBadRequest, "validation failed")
		default:
			logger.Log.Error("create news failed", "error", err)
			utils.WriteError(w, http.StatusInternalServerError, "failed to create news")
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, n)
}

// UpdateNews godoc
// @Summary      Update news
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id     path   int          true  "News ID"
// @Param        input  body   news.News  true  "News input"
// @Success      200  {object}  models.News
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Security     BearerAuth
// @Router       /api/news/{id} [put]
func (h *NewsHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var n models.News
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	n.ID = id

	actor, ok := getActor(r)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.newsService.UpdateNews(r.Context(), actor, &n); err != nil {
		switch {
		case errors.Is(err, errors2.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, "forbidden")
		case errors.Is(err, errors2.ErrNotFound):
			utils.WriteError(w, http.StatusNotFound, "news not found")
		case errors.Is(err, errors2.ErrValidation):
			utils.WriteError(w, http.StatusBadRequest, "validation failed")
		default:
			logger.Log.Error("update news failed", "error", err)
			utils.WriteError(w, http.StatusInternalServerError, "failed to update news")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, n)
}

// DeleteNews godoc
// @Summary      Delete news
// @Tags         news
// @Produce      json
// @Param        id   path   int  true  "News ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Security     BearerAuth
// @Router       /api/news/{id} [delete]
func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	actor, ok := getActor(r)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.newsService.DeleteNews(r.Context(), actor, id); err != nil {
		switch {
		case errors.Is(err, errors2.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, "forbidden")
		case errors.Is(err, errors2.ErrNotFound):
			utils.WriteError(w, http.StatusNotFound, "news not found")
		case errors.Is(err, errors2.ErrValidation):
			utils.WriteError(w, http.StatusBadRequest, "validation failed")
		default:
			logger.Log.Error("delete news failed", "error", err)
			utils.WriteError(w, http.StatusInternalServerError, "failed to delete news")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func getActor(r *http.Request) (models.Actor, bool) {
	actor, ok := r.Context().Value(middleware.CtxActor).(models.Actor)
	return actor, ok
}
