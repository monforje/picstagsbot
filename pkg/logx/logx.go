package logx

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	// handler := slog.NewTextHandler(os.Stdout, opts)
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
}

func Info(format string, args ...any) {
	logger.Info(format, args...)
}

func Error(format string, args ...any) {
	logger.Error(format, args...)
}

func Warn(format string, args ...any) {
	logger.Warn(format, args...)
}

func Fatal(format string, args ...any) {
	logger.Error(format, args...)
	os.Exit(1)
}
