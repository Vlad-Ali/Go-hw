package main

import (
	"log/slog"
	"transfersservice/internal/app"
)

func main() {
	cfg := &app.Config{Addr: ":8080", SecretKey: "secret"}
	application := app.NewApp(cfg)
	if err := application.Start(); err != nil {
		slog.Error("Error starting application", "error", err)
		return
	}
	slog.Info("Application stopped successfully")
}
