package balanceservice

import (
	"context"
	"transfersservice/internal/domain/balance"
	balanceobject "transfersservice/internal/domain/balance/object"
	"transfersservice/internal/domain/user/object"
)

type Service interface {
	GetBalance(ctx context.Context, userID object.UserID) (*balance.Balance, error)
	CreateBalance(ctx context.Context, profile balanceobject.BalanceProfile) (*balance.Balance, error)
	Deposit(ctx context.Context, userID object.UserID, amount balanceobject.Money) (*balance.Balance, error)
	Transfer(ctx context.Context, senderUserID object.UserID, targetEmail string, amount balanceobject.Money) (*balance.Balance, error)
}
