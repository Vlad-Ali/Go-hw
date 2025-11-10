package userrepo

import (
	"context"
	"transfersservice/internal/domain/user"
	"transfersservice/internal/domain/user/object"
)

type Repository interface {
	GetByUserID(ctx context.Context, id object.UserID) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, user *user.User) (*user.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
