package app

import (
	"net/http"
	"transfersservice/internal/adapter/balance"
	"transfersservice/internal/adapter/middleware"
	"transfersservice/internal/adapter/user"
	balanceservice "transfersservice/internal/domain/balance/service"
	userservice "transfersservice/internal/domain/user/service"
	"transfersservice/internal/domain/user/tokenservice"
)

type Handlers struct {
	UserHandler    *user.UserHandler
	BalanceHandler *balance.BalanceHandler
	AuthHandler    *middleware.AuthMiddleware
}

func NewHandlers(userService userservice.Service, balanceService balanceservice.Service, tokenService tokenservice.Service) *Handlers {
	userHandler := user.NewUserHandler(userService)
	balanceHandler := balance.NewBalanceHandler(balanceService)
	authHandler := middleware.NewAuthMiddleware(tokenService)
	return &Handlers{userHandler, balanceHandler, authHandler}
}

func (Handlers *Handlers) registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/user/register", Handlers.UserHandler.Register)
	mux.HandleFunc("POST /api/user/auth", Handlers.UserHandler.Authenticate)
	mux.HandleFunc("GET /api/user", Handlers.UserHandler.GetUser)

	mux.HandleFunc("POST /api/balance", Handlers.BalanceHandler.Create)
	mux.HandleFunc("GET /api/balance", Handlers.BalanceHandler.GetBalance)
	mux.HandleFunc("PUT /api/balance/deposit", Handlers.BalanceHandler.Deposit)
	mux.HandleFunc("PUT /api/balance/transfer", Handlers.BalanceHandler.Transfer)

	mainHandler := Handlers.AuthHandler.Authorize(mux)
	return mainHandler
}
