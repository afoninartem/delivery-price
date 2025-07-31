package l

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func InitLogger() {
	file, err := os.OpenFile("logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error(fmt.Sprintln(err))
		os.Exit(1)
	}
	defer file.Close()

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		slog.Error("can't get local time", "error", err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().In(loc).Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(file, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	slog.Info("start logging")
}
