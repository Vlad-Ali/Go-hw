package object

import (
	balanceerror "transfersservice/internal/domain/balance/error"
)

type Money struct {
	amount int
}

func NewMoney(amount int) (Money, error) {
	if amount < 0 {
		return Money{}, balanceerror.ErrMoneyCanNotBeNegative
	}
	return Money{amount: amount}, nil
}

func (m Money) WithAdd(amount int) (Money, error) {
	if m.amount+amount < 0 {
		return Money{}, balanceerror.ErrBalanceIsNotEnough
	}
	return Money{m.amount + amount}, nil
}

func (m Money) Amount() int {
	return m.amount
}
