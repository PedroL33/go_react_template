package util

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {
	Log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	Log.Info("Logger initialized.")
}
