package database

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"log/slog"
	"news-api/internal/config"
	"news-api/pkg/logger"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dbConfig := config.LoadConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Database.Host, dbConfig.Database.Port, dbConfig.Database.User, dbConfig.Database.Password, dbConfig.Database.DBName, dbConfig.Database.SSLMode,
	)

	logger.Log.Info("Connecting to News-api database...",
		slog.String("host", dbConfig.Database.Host),
		slog.String("dbname", dbConfig.Database.DBName),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Failed to open News-api database connection",
			slog.String("error", err.Error()))
		panic(err)
	}

	for i := 0; i < 30; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		logger.Log.Warn("Database not ready, retrying in 1s...",
			slog.String("error", err.Error()))
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		logger.Log.Error("Failed to ping News-api database after retries",
			slog.String("error", err.Error()))
		panic(err)
	}

	logger.Log.Info("Connected to News-api database successfully")

	if err := goose.Up(DB, "migrations"); err != nil {
		logger.Log.Error("Failed to run migrations",
			slog.String("error", err.Error()))
		panic(err)
	}
	logger.Log.Info("Migrations applied successfully")
}
