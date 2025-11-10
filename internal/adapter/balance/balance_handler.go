package balance

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"transfersservice/internal/adapter/balance/request"
	balanceresponse "transfersservice/internal/adapter/balance/response"
	"transfersservice/internal/adapter/useridkey"
	balanceerror "transfersservice/internal/domain/balance/error"
	"transfersservice/internal/domain/balance/object"
	balanceservice "transfersservice/internal/domain/balance/service"
	usererror "transfersservice/internal/domain/user/error"
)

type BalanceHandler struct {
	balanceService balanceservice.Service
}

func NewBalanceHandler(balanceService balanceservice.Service) *BalanceHandler {
	return &BalanceHandler{balanceService}
}

func (b *BalanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Debug("BalanceHandler.Create called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error extracting user id", "error", err)
		http.Error(w, "Failed to authorize user", http.StatusUnauthorized)
		return
	}

	var balanceCreateRequest request.BalanceCreateRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error reading body", "error", err)
		http.Error(w, "Failed to Create balancer", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &balanceCreateRequest)
	if err != nil {
		slog.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	balance, err := b.balanceService.CreateBalance(r.Context(), object.NewBalanceProfile(userID, balanceCreateRequest.Name))

	if err != nil {
		slog.Error("Error creating balance", "error", err)
		if errors.Is(err, balanceerror.ErrBalanceNameValidationFailed) {
			http.Error(w, "Name is empty", http.StatusBadRequest)
		} else if errors.Is(err, balanceerror.ErrBalanceAlreadyExists) {
			http.Error(w, "Balance already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create balance", http.StatusInternalServerError)
		}
		return
	}

	response := balanceresponse.BalanceCreateResponse{Name: balance.Name(), Amount: balance.Money().Amount()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}

func (b *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	slog.Debug("BalanceHandler.GetBalance called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error extracting user id", "error", err)
		http.Error(w, "Failed to authorize user", http.StatusUnauthorized)
		return
	}

	balance, err := b.balanceService.GetBalance(r.Context(), userID)
	if err != nil {
		slog.Error("Error getting balance", "error", err)
		if errors.Is(err, balanceerror.ErrBalanceIsNotFound) {
			http.Error(w, "Balance not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		}
		return
	}

	response := balanceresponse.BalanceGetResponse{Name: balance.Name(), Balance: balance.Money().Amount()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}

func (b *BalanceHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	slog.Debug("BalanceHandler.Deposit called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error extracting user id", "error", err)
		http.Error(w, "Failed to authorize user", http.StatusUnauthorized)
		return
	}

	var balanceDepositRequest request.BalanceDepositRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error reading body", "error", err)
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &balanceDepositRequest)
	if err != nil {
		slog.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	money, err := object.NewMoney(balanceDepositRequest.Amount)
	if err != nil {
		slog.Error("Error creating money", "error", err)
		http.Error(w, "Amount cannot be negative", http.StatusBadRequest)
		return
	}

	balance, err := b.balanceService.Deposit(r.Context(), userID, money)
	if err != nil {
		slog.Error("Error depositing balance", "error", err)
		if errors.Is(err, balanceerror.ErrBalanceIsNotFound) {
			http.Error(w, "Balance not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to deposit balance", http.StatusInternalServerError)
		}
		return
	}

	response := balanceresponse.BalanceDepositResponse{Name: balance.Name(), Amount: balance.Money().Amount()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}

func (b *BalanceHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	slog.Debug("BalanceHandler.Transfer called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error extracting user id", "error", err)
		http.Error(w, "Failed to authorize user", http.StatusUnauthorized)
		return
	}

	var balanceTransferRequest request.BalanceTransferRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error reading body", "error", err)
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &balanceTransferRequest)
	if err != nil {
		slog.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	money, err := object.NewMoney(balanceTransferRequest.Amount)
	if err != nil {
		slog.Error("Error creating money", "error", err)
		http.Error(w, "Amount cannot be negative", http.StatusBadRequest)
		return
	}

	balance, err := b.balanceService.Transfer(r.Context(), userID, balanceTransferRequest.TargetEmail, money)
	if err != nil {
		slog.Error("Error transfering balance", "error", err)
		if errors.Is(err, balanceerror.ErrTransferBalanceToYourself) {
			http.Error(w, "It is forbidden to transfer to yourself", http.StatusBadRequest)
		} else if errors.Is(err, usererror.ErrUserIsNotFound) || errors.Is(err, balanceerror.ErrBalanceIsNotFound) {
			http.Error(w, "Target balance is not found", http.StatusNotFound)
		} else if errors.Is(err, balanceerror.ErrBalanceIsNotEnough) {
			http.Error(w, "Money is not enough on balance", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to transfer money", http.StatusInternalServerError)
		}
		return
	}

	response := balanceresponse.BalanceTransferResponse{Balance: balance.Money().Amount(), Name: balance.Name()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}
