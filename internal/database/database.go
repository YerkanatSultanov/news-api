package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"news-api/internal/config"
	"news-api/pkg/logger"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	logger.Log.Info("Connecting to News-api database...",
		slog.String("host", dbConfig.Host),
		slog.String("dbname", dbConfig.DBName),
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
}
