package balancevalidation

import (
	balanceerror "transfersservice/internal/domain/balance/error"
)

func ValidateBalanceName(name string) error {
	if name == "" {
		return balanceerror.ErrBalanceNameValidationFailed
	}
	return nil
}
