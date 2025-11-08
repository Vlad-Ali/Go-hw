package object

import (
	balanceerror "transfersservice/internal/domain/balance/error"

	"github.com/google/uuid"
)

type BalanceId struct {
	id string
}

func NewBalanceId(id string) (BalanceId, error) {
	if _, err := uuid.Parse(id); err != nil {
		return BalanceId{}, balanceerror.ErrBalanceIDCreatingIsNotValid
	}
	return BalanceId{id: id}, nil
}

func (b BalanceId) ID() string {
	return b.id
}

func (b BalanceId) IsEmpty() bool {
	return b.id == ""
}
