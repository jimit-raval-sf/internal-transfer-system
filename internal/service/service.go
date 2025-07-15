package service

import (
	"errors"
	"fmt"
	"log"

	"internal-transfer-system/internal/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

type CreateAccountRequest struct {
	AccountID      uint   `json:"account_id" binding:"required"`
	InitialBalance string `json:"initial_balance" binding:"required"`
}

type AccountResponse struct {
	AccountID uint   `json:"account_id"`
	Balance   string `json:"balance"`
}

type CreateTransactionRequest struct {
	SourceAccountID      uint   `json:"source_account_id" binding:"required"`
	DestinationAccountID uint   `json:"destination_account_id" binding:"required"`
	Amount               string `json:"amount" binding:"required"`
}

func (s *Service) CreateAccount(req *CreateAccountRequest) error {
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

	if s.repo.AccountExists(req.AccountID) {
		return errors.New("account already exists")
	}

	account := &repository.Account{
		AccountID: req.AccountID,
		Balance:   balance,
	}

	err = s.repo.CreateAccount(account)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		return errors.New("failed to create account")
	}

	return nil
}

func (s *Service) GetAccount(accountID uint) (*AccountResponse, error) {
	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		log.Printf("Error getting account: %v", err)
		return nil, errors.New("failed to get account")
	}

	return &AccountResponse{
		AccountID: account.AccountID,
		Balance:   account.Balance.String(),
	}, nil
}

func (s *Service) CreateTransaction(req *CreateTransactionRequest) error {
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

	if !s.repo.AccountExists(req.SourceAccountID) {
		return fmt.Errorf("source account %d not found", req.SourceAccountID)
	}

	if !s.repo.AccountExists(req.DestinationAccountID) {
		return fmt.Errorf("destination account %d not found", req.DestinationAccountID)
	}

	err = s.repo.CreateTransaction(req.SourceAccountID, req.DestinationAccountID, amount)
	if err != nil {
		if err.Error() == "insufficient balance" {
			return errors.New("insufficient balance")
		}
		log.Printf("Error creating transaction: %v", err)
		return errors.New("failed to create transaction")
	}

	return nil
}