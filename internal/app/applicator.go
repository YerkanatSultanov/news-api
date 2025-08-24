package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"news-api/internal/config"
	"news-api/internal/database"
	"news-api/internal/http/handlers"
	"news-api/internal/http/router"
	"news-api/internal/repository"
	"news-api/internal/service"
	"news-api/pkg/logger"
	redisClient "news-api/pkg/redis"
	"news-api/pkg/token"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	DB          *sql.DB
	AuthRepo    *repository.UserRepository
	AuthService *service.AuthService
	AuthHandler *handlers.AuthHandler
	JWTManager  *token.JWTManager
	RedisClient *redis.Client
	server      *http.Server
}

func NewApp() *App {
	config.LoadConfig()
	cfg := config.AppConfig

	logger.InitLogger(cfg.Log.Level)
	logger.Log.Info("Logger initialized", "level", cfg.Log.Level)

	logger.Log.Info("Initializing database...")
	database.InitDB()
	logger.Log.Info("Database initialized")

	logger.Log.Info("Initializing Redis client...")
	client := redisClient.InitRedis(cfg.Redis)
	if client == nil {
		panic("Failed to initialize Redis")
	}

	jwtManager := token.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	authRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(authRepo, client, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	return &App{
		DB:          database.DB,
		AuthRepo:    authRepo,
		AuthService: authService,
		AuthHandler: authHandler,
		JWTManager:  jwtManager,
		RedisClient: client,
	}
}

func (a *App) Run() {
	cfg := config.AppConfig

	routers := router.NewRouter(a.AuthHandler)
	a.server = &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: routers,
	}

	go func() {
		logger.Log.Info("HTTP server started", "port", cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Shutdown(ctx); err != nil {
		logger.Log.Error("Shutdown finished with errors", "error", err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	logger.Log.Info("Shutting down server...")

	var errList []error

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Log.Error("Error shutting down HTTP server", "error", err)
			errList = append(errList, err)
		}
	}

	if a.RedisClient != nil {
		if err := a.RedisClient.Close(); err != nil {
			logger.Log.Error("Error closing Redis", "error", err)
			errList = append(errList, err)
		}
	}

	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			logger.Log.Error("Error closing DB", "error", err)
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("shutdown errors: %v", errList)
	}

	logger.Log.Info("Shutdown complete")
	return nil
}
