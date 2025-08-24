package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"news-api/internal/dto/auth"
	"news-api/internal/service/interfaces"
	"news-api/pkg/logger"
	"news-api/utils"
)

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input auth.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.Warn("Invalid register request", slog.String("error", err.Error()))
		utils.WriteJSON(w, http.StatusBadRequest, auth.Response{Message: "Invalid request payload"})
		return
	}

	if err := h.authService.Register(r.Context(), input); err != nil {
		logger.Log.Warn("User registration failed", slog.String("email", input.Email), slog.String("error", err.Error()))
		utils.WriteJSON(w, http.StatusBadRequest, auth.Response{Message: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, auth.Response{Message: "User successfully registered"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input auth.LoginUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.Warn("Invalid login request", slog.String("error", err.Error()))
		utils.WriteJSON(w, http.StatusBadRequest, auth.Response{Message: "Invalid request payload"})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), input)
	if err != nil || accessToken == "" {
		logger.Log.Warn("Login failed", slog.String("email", input.Email))
		utils.WriteJSON(w, http.StatusUnauthorized, auth.Response{Message: "Invalid email or password"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, auth.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
