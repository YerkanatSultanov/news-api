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

// Register godoc
// @Summary      Register new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body   auth.RegisterUserInput  true  "User registration data"
// @Success      201  {object}  auth.Response
// @Failure      400  {object}  auth.Response
// @Router       /api/register [post]
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

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body   auth.LoginUserInput  true  "Login credentials"
// @Success      200  {object}  auth.TokensResponse
// @Failure      400  {object}  auth.Response
// @Failure      401  {object}  auth.Response
// @Router       /api/login [post]
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

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate refresh token for the current user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  auth.Response
// @Failure      401  {object}  auth.Response
// @Failure      500  {object}  auth.Response
// @Security     BearerAuth
// @Router       /api/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	actor, ok := getActor(r)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authService.Logout(r.Context(), actor.UserID); err != nil {
		logger.Log.Error("Logout failed", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	utils.WriteJSON(w, http.StatusOK, auth.Response{Message: "Logout successfully"})
}
