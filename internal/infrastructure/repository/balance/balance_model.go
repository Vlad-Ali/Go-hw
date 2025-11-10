package balancereppository

import (
	balancedomain "transfersservice/internal/domain/balance"
	"transfersservice/internal/domain/balance/object"
	userobject "transfersservice/internal/domain/user/object"
)

type BalanceModel struct {
	ID     string
	Name   string
	UserID string
	Amount int
}

func NewBalanceModel(id string, name string, userID string, amount int) *BalanceModel {
	return &BalanceModel{id, name, userID, amount}
}

func (b *BalanceModel) ToDomain() (*balancedomain.Balance, error) {
	userID, err := userobject.NewUserID(b.UserID)
	if err != nil {
		return nil, err
	}

	balanceID, err := object.NewBalanceId(b.ID)
	if err != nil {
		return nil, err
	}

	balance := balancedomain.NewBalance(b.Name, userID)
	err = balance.AddMoney(b.Amount)
	if err != nil {
		return nil, err
	}

	err = balance.SetID(balanceID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
