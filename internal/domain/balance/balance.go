package balance

import (
	balanceerror "transfersservice/internal/domain/balance/error"
	balanceobject "transfersservice/internal/domain/balance/object"
	"transfersservice/internal/domain/user/object"
)

type Balance struct {
	name   string
	money  balanceobject.Money
	userID object.UserID
	id     balanceobject.BalanceId
}

func NewBalance(name string, userID object.UserID) *Balance {
	return &Balance{name: name, userID: userID}
}

func (b *Balance) Name() string {
	return b.name
}

func (b *Balance) Money() balanceobject.Money {
	return b.money
}

func (b *Balance) UserID() object.UserID {
	return b.userID
}

func (b *Balance) ID() balanceobject.BalanceId {
	return b.id
}

func (b *Balance) AddMoney(amount int) error {
	money, err := b.money.WithAdd(amount)
	if err != nil {
		return err
	}
	b.money = money
	return nil
}

func (b *Balance) SetID(id balanceobject.BalanceId) error {
	if b.id.IsEmpty() {
		b.id = id
		return nil
	}
	return balanceerror.ErrBalanceIDAlreadyExists
}

func (b *Balance) Transfer(target *Balance, amount int) error {
	money, err := b.money.WithAdd(-amount)
	if err != nil {
		return err
	}

	targetMoney, err := target.money.WithAdd(amount)
	if err != nil {
		return err
	}

	b.money = money
	target.money = targetMoney
	return nil
}
