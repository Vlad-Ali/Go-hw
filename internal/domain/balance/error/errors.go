package error

import "errors"

var (
	ErrBalanceIDCreatingIsNotValid = errors.New("balanceID is not valid")
	ErrBalanceIsNotFound           = errors.New("balance is not found")
	ErrBalanceIsNotEnough          = errors.New("balance is not enough")
	ErrMoneyCanNotBeNegative       = errors.New("money can not be negative")
	ErrBalanceAlreadyExists        = errors.New("balance already exists")
	ErrBalanceIDAlreadyExists      = errors.New("balance id already exists")
	ErrBalanceNameValidationFailed = errors.New("failed to validate balance name")
)
