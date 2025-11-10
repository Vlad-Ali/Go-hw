package app

import (
	balancerepo "transfersservice/internal/domain/balance/repository"
	userrepo "transfersservice/internal/domain/user/repository"
	balancereppository "transfersservice/internal/infrastructure/repository/balance"
	userrepository "transfersservice/internal/infrastructure/repository/user"
)

type Repositories struct {
	BalanceRepo balancerepo.Repository
	UserRepo    userrepo.Repository
}

func NewRepositories(balanceStorage map[string]*balancereppository.BalanceModel, userStorage map[string]*userrepository.UserModel) *Repositories {
	return &Repositories{balancereppository.NewBalanceRepository(balanceStorage), userrepository.NewUserRepository(userStorage)}
}
