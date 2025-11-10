package balanceservice

import (
	"context"
	"fmt"
	"log/slog"
	balancevalidation "transfersservice/internal/application/usecase/balance/validation"
	balancedomain "transfersservice/internal/domain/balance"
	balanceerror "transfersservice/internal/domain/balance/error"
	balanceobject "transfersservice/internal/domain/balance/object"
	balancerepo "transfersservice/internal/domain/balance/repository"
	"transfersservice/internal/domain/user/object"
	userrepo "transfersservice/internal/domain/user/repository"
)

type BalanceService struct {
	balanceRepo balancerepo.Repository
	userRepo    userrepo.Repository
}

func NewBalanceService(balanceRepo balancerepo.Repository, userRepo userrepo.Repository) *BalanceService {
	return &BalanceService{balanceRepo: balanceRepo, userRepo: userRepo}
}

func (b *BalanceService) GetBalance(ctx context.Context, userID object.UserID) (*balancedomain.Balance, error) {
	balance, err := b.balanceRepo.GetByUserID(ctx, userID.ID())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get balance for user with id: %v error: %v", userID, err))
		return nil, err
	}
	slog.Debug(fmt.Sprintf("get balance for user with id: %v", userID))
	return balance, nil
}

func (b *BalanceService) CreateBalance(ctx context.Context, profile balanceobject.BalanceProfile) (*balancedomain.Balance, error) {
	err := balancevalidation.ValidateBalanceName(profile.Name())
	if err != nil {
		slog.Error("failed to validate balance name", "error", err)
		return nil, err
	}

	ok, err := b.balanceRepo.ExistsByUserID(ctx, profile.UserID().ID())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to check existence for balance for user with id: %v error: %v", profile.UserID(), err))
		return nil, err
	}
	if ok {
		slog.Error(fmt.Sprintf("balance already exists for user with id: %v", profile.UserID()))
		return nil, balanceerror.ErrBalanceAlreadyExists
	}

	balance := balancedomain.NewBalance(profile.Name(), profile.UserID())
	balance, err = b.balanceRepo.Save(ctx, balance)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create balance for user with id: %v error: %v", profile.UserID(), err))
		return nil, err
	}
	slog.Debug(fmt.Sprintf("create balance for user with id: %v", profile.UserID()))
	return balance, nil
}

func (b *BalanceService) Deposit(ctx context.Context, userID object.UserID, amount balanceobject.Money) (*balancedomain.Balance, error) {
	balance, err := b.GetBalance(ctx, userID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to deposit balance with user id: %v error: %v", userID, err))
		return nil, err
	}

	err = balance.AddMoney(amount.Amount())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to deposit balance with user id: %v error: %v", userID, err))
		return nil, err
	}

	balance, err = b.balanceRepo.Save(ctx, balance)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to deposit balance for user with id: %v error: %v", userID, err))
		return nil, err
	}

	slog.Debug("deposit balance for user with id", "userID", userID.ID())
	return balance, nil
}

func (b *BalanceService) Transfer(ctx context.Context, senderUserID object.UserID, targetEmail string, amount balanceobject.Money) (*balancedomain.Balance, error) {
	senderBalance, err := b.balanceRepo.GetByUserID(ctx, senderUserID.ID())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to find sender balance with user id: %v error: %v", senderUserID, err))
		return nil, err
	}

	targetUser, err := b.userRepo.GetByEmail(ctx, targetEmail)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to find user with email: %v error: %v", targetEmail, err))
		return nil, err
	}

	if targetUser.ID().ID() == senderUserID.ID() {
		slog.Error("failed to transfer to yourself", "error", err)
		return nil, balanceerror.ErrTransferBalanceToYourself
	}

	targetBalance, err := b.balanceRepo.GetByUserID(ctx, targetUser.ID().ID())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to find target balance with user id: %v error: %v", targetUser.ID(), err))
		return nil, err
	}

	err = senderBalance.Transfer(targetBalance, amount.Amount())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to transfer money sender id: %v error: %v", senderUserID, err))
		return nil, err
	}

	senderBalance, err = b.balanceRepo.Save(ctx, senderBalance)
	if err != nil {
		slog.Error("failed to transfer sender balance", "error", err)
		return nil, err
	}

	targetBalance, err = b.balanceRepo.Save(ctx, targetBalance)
	if err != nil {
		slog.Error("failed to transfer target balance", "error", err)
		return nil, err
	}

	slog.Debug("successfully transfer balance for user with id", "targetID", targetUser.ID().ID())
	return senderBalance, nil
}
