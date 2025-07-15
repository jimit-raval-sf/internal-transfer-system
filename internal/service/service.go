package service

import (
	"errors"
	"log"

	"internal-transfer-system/internal/repository"
	"internal-transfer-system/internal/validator"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}


type AccountResponse struct {
	AccountID uint   `json:"account_id"`
	Balance   string `json:"balance"`
}


func (s *Service) CreateAccount(req *validator.CreateAccountRequest) error {
	if err := validator.ValidateCreateAccountRequest(req, s.repo); err != nil {
		return err
	}

	balance, _ := decimal.NewFromString(req.InitialBalance)

	account := &repository.Account{
		AccountID: req.AccountID,
		Balance:   balance,
	}

	err := s.repo.CreateAccount(account)
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

func (s *Service) CreateTransaction(req *validator.CreateTransactionRequest) error {
	if err := validator.ValidateCreateTransactionRequest(req, s.repo); err != nil {
		return err
	}

	amount, _ := decimal.NewFromString(req.Amount)

	err := s.repo.CreateTransaction(req.SourceAccountID, req.DestinationAccountID, amount)
	if err != nil {
		if err.Error() == "insufficient balance" {
			return errors.New("insufficient balance")
		}
		log.Printf("Error creating transaction: %v", err)
		return errors.New("failed to create transaction")
	}

	return nil
}