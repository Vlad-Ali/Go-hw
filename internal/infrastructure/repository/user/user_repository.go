package userrepository

import (
	"context"
	"log/slog"
	userdomain "transfersservice/internal/domain/user"
	usererror "transfersservice/internal/domain/user/error"
	"transfersservice/internal/domain/user/object"

	"github.com/google/uuid"
)

type UserRepository struct {
	mapStorage map[string]*UserModel
}

func NewUserRepository(mapStorage map[string]*UserModel) *UserRepository {
	return &UserRepository{mapStorage: mapStorage}
}

func (u *UserRepository) GetByUserID(ctx context.Context, id object.UserID) (*userdomain.User, error) {
	userModel, ok := u.mapStorage[id.ID()]
	if !ok {
		slog.Error("user not found with id", "ID", id.ID())
		return nil, usererror.ErrUserIsNotFound
	}

	user, err := userModel.ToDomain()
	if err != nil {
		slog.Error("failed to create userDomain with id", "ID", id.ID())
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	for _, userModel := range u.mapStorage {
		if userModel.Email == email {
			user, err := userModel.ToDomain()
			if err != nil {
				slog.Error("failed to create userDomain with id", "ID", userModel.ID)
				return nil, err
			}
			return user, nil
		}
	}
	slog.Error("user not found with email", "email", email)
	return nil, usererror.ErrUserIsNotFound
}

func (u *UserRepository) Create(ctx context.Context, user *userdomain.User) (*userdomain.User, error) {
	userModel := NewUserModel(user.Username(), user.Password(), user.Email(), uuid.New().String())
	u.mapStorage[userModel.ID] = userModel
	slog.Debug("user created with id", "ID", userModel.ID)
	return userModel.ToDomain()
}

func (u *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	for _, userModel := range u.mapStorage {
		if userModel.Email == email {
			return true, nil
		}
	}
	return false, nil
}
