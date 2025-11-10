package userservice

import (
	"context"
	"transfersservice/internal/domain/user"
	"transfersservice/internal/domain/user/object"
)

type Service interface {
	GetUserByID(ctx context.Context, id object.UserID) (*user.User, error)
	Register(ctx context.Context, data object.UserRegistrationData) (*user.User, error)
	Authenticate(ctx context.Context, data object.AuthenticationData) (string, error)
}
