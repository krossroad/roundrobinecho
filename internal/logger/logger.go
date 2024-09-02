package logger

import (
	"log/slog"
	"os"
)

func New(app string) *slog.Logger {
	const version = "1.0.0-beta1"

	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(h)
	return logger.With(
		slog.Group(
			"app-info",
			slog.String("version", version),
			slog.String("server-id", os.Getenv("HOSTNAME")),
			slog.String("app", app),
		),
	)
}
