package validator

import (
	"errors"
	"fmt"

	"internal-transfer-system/internal/repository"

	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	AccountID      uint   `json:"account_id" binding:"required"`
	InitialBalance string `json:"initial_balance" binding:"required"`
}

type CreateTransactionRequest struct {
	SourceAccountID      uint   `json:"source_account_id" binding:"required"`
	DestinationAccountID uint   `json:"destination_account_id" binding:"required"`
	Amount               string `json:"amount" binding:"required"`
}

func ValidateCreateAccountRequest(req *CreateAccountRequest, repo *repository.Repository) error {
	if req.AccountID == 0 {
		return errors.New("account_id must be a positive integer")
	}

	balance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return errors.New("invalid initial_balance format")
	}

	if balance.IsNegative() {
		return errors.New("initial_balance must be non-negative")
	}

	if balance.Exponent() < -5 {
		return errors.New("initial_balance must have at most 5 decimal places")
	}

	if repo.AccountExists(req.AccountID) {
		return errors.New("account already exists")
	}

	return nil
}

func ValidateCreateTransactionRequest(req *CreateTransactionRequest, repo *repository.Repository) error {
	if req.SourceAccountID == req.DestinationAccountID {
		return errors.New("cannot transfer to the same account")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return errors.New("invalid amount format")
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than 0")
	}

	if amount.Exponent() < -5 {
		return errors.New("amount must have at most 5 decimal places")
	}

	if !repo.AccountExists(req.SourceAccountID) {
		return fmt.Errorf("source account %d not found", req.SourceAccountID)
	}

	if !repo.AccountExists(req.DestinationAccountID) {
		return fmt.Errorf("destination account %d not found", req.DestinationAccountID)
	}

	return nil
}