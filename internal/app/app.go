package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	balancereppository "transfersservice/internal/infrastructure/repository/balance"
	userrepository "transfersservice/internal/infrastructure/repository/user"
)

type App struct {
	server       *http.Server
	repositories *Repositories
	services     *Services
	handlers     *Handlers
}

type Config struct {
	Addr      string
	SecretKey string
}

func NewApp(cfg *Config) *App {
	repositories := NewRepositories(make(map[string]*balancereppository.BalanceModel), make(map[string]*userrepository.UserModel))
	services := NewServices(cfg.SecretKey, repositories.BalanceRepo, repositories.UserRepo)
	handlers := NewHandlers(services.UserService, services.BalanceService, services.TokenService)

	handler := handlers.registerRoutes()

	setLogger()
	return &App{
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: handler,
		}, repositories: repositories, services: services,
		handlers: handlers,
	}
}

func (app *App) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		slog.Info(fmt.Sprintf("Server started at %s", app.server.Addr))
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("Server failed: %v", err))
			os.Exit(1)
		}
	}()

	<-signalChan
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return app.server.Shutdown(ctx)
}

func setLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
