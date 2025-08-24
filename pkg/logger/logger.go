package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

var Log *slog.Logger

func InitLogger(level string) {
	var slogLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn", "warning":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	logHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slogLevel,
		TimeFormat: "15:04:05",
	})
	Log = slog.New(logHandler)
}
