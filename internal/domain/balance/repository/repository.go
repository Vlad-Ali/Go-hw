package balancerepository

import (
	"context"
	"transfersservice/internal/domain/balance"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*balance.Balance, error)
	Save(ctx context.Context, balance *balance.Balance) (*balance.Balance, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}
