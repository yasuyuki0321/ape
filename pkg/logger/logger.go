package logger

import (
	"log"
	"os"

	"golang.org/x/exp/slog"

	"github.com/yasuyuki0321/ape/pkg/utils"
)

var logger *slog.Logger

func init() {
	logFile, err := os.OpenFile(utils.GetHomePath("~/.psh_hisotry"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	ops := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger = slog.New(slog.NewJSONHandler(logFile, &ops))
	slog.SetDefault(logger)
}
