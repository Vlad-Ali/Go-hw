package balancereppository

import (
	"context"
	"log/slog"
	balancedomain "transfersservice/internal/domain/balance"
	balanceerror "transfersservice/internal/domain/balance/error"
	balanceobject "transfersservice/internal/domain/balance/object"

	"github.com/google/uuid"
)

type BalanceRepository struct {
	mapStorage map[string]*BalanceModel
}

func NewBalanceRepository(mapStorage map[string]*BalanceModel) *BalanceRepository {
	return &BalanceRepository{mapStorage}
}

func (b *BalanceRepository) GetByUserID(ctx context.Context, userID string) (*balancedomain.Balance, error) {
	for _, balanceModel := range b.mapStorage {
		if balanceModel.UserID == userID {
			balance, err := balanceModel.ToDomain()
			if err != nil {
				slog.Error("fail to convert to domain model", "error", err)
				return nil, err
			}
			return balance, nil
		}
	}
	return nil, balanceerror.ErrBalanceIsNotFound
}

func (b *BalanceRepository) Save(ctx context.Context, balance *balancedomain.Balance) (*balancedomain.Balance, error) {
	_, ok := b.mapStorage[balance.ID().ID()]
	if ok {
		model := b.mapStorage[balance.ID().ID()]
		model.Amount = balance.Money().Amount()
		b.mapStorage[balance.ID().ID()] = model
		return balance, nil
	}
	balanceID, err := balanceobject.NewBalanceId(uuid.New().String())
	if err != nil {
		slog.Error("fail to generate balance id", "error", err)
		return nil, err
	}
	balanceModel := NewBalanceModel(balanceID.ID(), balance.Name(), balance.UserID().ID(), balance.Money().Amount())
	b.mapStorage[balanceID.ID()] = balanceModel
	return balanceModel.ToDomain()
}

func (b *BalanceRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	for _, balanceModel := range b.mapStorage {
		if balanceModel.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}
