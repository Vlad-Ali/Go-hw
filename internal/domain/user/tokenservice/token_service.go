package tokenservice

import (
	"context"
	"transfersservice/internal/domain/user"
	"transfersservice/internal/domain/user/object"
)

type Service interface {
	GenerateToken(ctx context.Context, user *user.User) (string, error)
	ValidateToken(ctx context.Context, token string) (object.UserID, error)
}
