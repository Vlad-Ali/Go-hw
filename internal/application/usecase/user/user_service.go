package userservice

import (
	"context"
	"fmt"
	"log/slog"
	"transfersservice/internal/application/usecase/user/hasher"
	uservalidation "transfersservice/internal/application/usecase/user/validation"
	userdomain "transfersservice/internal/domain/user"
	usererror "transfersservice/internal/domain/user/error"
	"transfersservice/internal/domain/user/object"
	userrepo "transfersservice/internal/domain/user/repository"
	"transfersservice/internal/domain/user/tokenservice"
)

type UserService struct {
	tokenService tokenservice.Service
	userRepo     userrepo.Repository
}

func NewUserService(tokenService tokenservice.Service, userRepo userrepo.Repository) *UserService {
	return &UserService{tokenService: tokenService, userRepo: userRepo}
}

func (s *UserService) Register(ctx context.Context, data object.UserRegistrationData) (*userdomain.User, error) {
	err := uservalidation.ValidateUserRegistrationData(data)
	if err != nil {
		slog.Error("Validation failed", "error", err)
		return nil, err
	}
	ok, err := s.userRepo.ExistsByEmail(ctx, data.Email())
	if err != nil {
		slog.Error("Failed to check if email exists", "error", err)
		return nil, err
	}

	if ok {
		slog.Error("Email already exists")
		return nil, usererror.ErrUserEmailAlreadyExists
	}

	hashPassword, err := hasher.HashPassword(data.Password())
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return nil, err
	}

	user := userdomain.NewUser(data.Username(), hashPassword, data.Email())
	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return nil, usererror.ErrFailedToRegisterUser
	}
	slog.Debug("User created with id", "userID", user.ID().ID())
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, data object.AuthenticationData) (string, error) {
	err := uservalidation.ValidateAuthenticationData(data)
	if err != nil {
		slog.Error("Validation failed", "error", err)
		return "", err
	}

	user, err := s.userRepo.GetByEmail(ctx, data.Email())
	if err != nil {
		slog.Error("failed to auth user error", "error", err)
		return "", err
	}

	ok := hasher.VerifyPassword(data.Password(), user.Password())
	if !ok {
		slog.Error("failed to auth user")
		return "", usererror.ErrInvalidPassword
	}

	token, err := s.tokenService.GenerateToken(ctx, user)
	if err != nil {
		slog.Error("failed to generate token error", "error", err)
		return "", err
	}
	slog.Debug("user is authenticated")
	return token, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id object.UserID) (*userdomain.User, error) {
	user, err := s.userRepo.GetByUserID(ctx, id)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to find user: %s, %v", id, err))
		return nil, err
	}
	slog.Debug("user found with id", "ID", id.ID())
	return user, nil
}
