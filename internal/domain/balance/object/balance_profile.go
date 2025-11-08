package object

import (
	"transfersservice/internal/domain/user/object"
)

type BalanceProfile struct {
	userID object.UserID
	name   string
}

func NewBalanceProfile(userID object.UserID, name string) BalanceProfile {
	return BalanceProfile{userID, name}
}

func (b BalanceProfile) UserID() object.UserID { return b.userID }

func (b BalanceProfile) Name() string { return b.name }
