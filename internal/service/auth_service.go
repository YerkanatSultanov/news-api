package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"news-api/internal/dto/auth"
	"news-api/internal/models"
	"news-api/internal/repository/interfaces"
	"news-api/pkg/logger"
	"news-api/pkg/password"
	"news-api/pkg/token"
	"strconv"
	"time"
)

const (
	passwordLength = 6
	contextTimeout = 5 * time.Second
)

type AuthService struct {
	authRepo   interfaces.UserRepository
	redis      *redis.Client
	jwtManager *token.JWTManager
}

func NewAuthService(repo interfaces.UserRepository, redis *redis.Client, jwtManager *token.JWTManager) *AuthService {
	return &AuthService{authRepo: repo, redis: redis, jwtManager: jwtManager}
}

func (s *AuthService) Register(ctx context.Context, input auth.RegisterUserInput) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if len(input.Password) < passwordLength {
		err := fmt.Errorf("password must be at least 6 characters")
		logger.Log.Error("Password validation failed: " + err.Error())
		return err
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		logger.Log.Error("Failed to hash password: " + err.Error())
		return err
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      "editor",
		Avatar:    input.Avatar,
		CreatedAt: time.Now(),
	}

	return s.authRepo.Create(ctx, &user)
}

func (s *AuthService) Login(ctx context.Context, input auth.LoginUserInput) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	user, err := s.authRepo.GetByEmail(input.Email)
	if err != nil {
		logger.Log.Warn("Login failed: user not found", slog.String("email", input.Email), slog.String("error", err.Error()))
		return "", "", err
	}

	if !password.CheckPassword(input.Password, user.Password) {
		logger.Log.Warn("Login failed: incorrect password", slog.String("email", input.Email))
		return "", "", err
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokens(user.ID, user.Role)
	if err != nil {
		logger.Log.Error("Failed to generate tokens: " + err.Error())
		return "", "", err
	}

	if err := s.SaveRefreshToken(user.ID, refreshToken); err != nil {
		logger.Log.Error("Failed to save refresh token: " + err.Error())
		return "", "", err
	}
	logger.Log.Info("User logged in", slog.String("email", input.Email))
	return accessToken, refreshToken, nil

}

func (s *AuthService) SaveRefreshToken(userID int, refreshToken string) error {
	err := s.redis.Set(context.Background(), s.getRefreshTokenKey(userID), refreshToken, 24*time.Hour*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) getRefreshTokenKey(userID int) string {
	return "refresh_token:" + strconv.Itoa(userID)
}
