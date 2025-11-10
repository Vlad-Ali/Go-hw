package app

import (
	balanceserviceimpl "transfersservice/internal/application/usecase/balance"
	"transfersservice/internal/application/usecase/jwt"
	userserviceimpl "transfersservice/internal/application/usecase/user"
	balancerepo "transfersservice/internal/domain/balance/repository"
	balanceservice "transfersservice/internal/domain/balance/service"
	userrepo "transfersservice/internal/domain/user/repository"
	userservice "transfersservice/internal/domain/user/service"
	"transfersservice/internal/domain/user/tokenservice"
)

type Services struct {
	BalanceService balanceservice.Service
	UserService    userservice.Service
	TokenService   tokenservice.Service
}

func NewServices(secretKey string, balanceRepo balancerepo.Repository, userRepo userrepo.Repository) *Services {
	tokenService := jwt.NewJwtService(secretKey)
	userService := userserviceimpl.NewUserService(tokenService, userRepo)
	balanceService := balanceserviceimpl.NewBalanceService(balanceRepo, userRepo)
	return &Services{balanceService, userService, tokenService}
}
