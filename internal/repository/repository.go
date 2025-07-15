package repository

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAccount(account *Account) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if account already exists by AccountID
		var existingAccount Account
		err := tx.Where("account_id = ?", account.AccountID).First(&existingAccount).Error
		if err == nil {
			return errors.New("account already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		
		// Create account
		return tx.Create(account).Error
	})
}

func (r *Repository) GetAccountByID(accountID uint) (*Account, error) {
	var account Account
	err := r.db.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *Repository) AccountExists(accountID uint) bool {
	var count int64
	r.db.Model(&Account{}).Where("account_id = ?", accountID).Count(&count)
	return count > 0
}

func (r *Repository) CreateTransaction(sourceAccountID, destAccountID uint, amount decimal.Decimal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var sourceAccount, destAccount Account
		
		// Get source account by AccountID
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("account_id = ?", sourceAccountID).First(&sourceAccount).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("source account not found")
			}
			return err
		}
		
		// Get destination account by AccountID
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("account_id = ?", destAccountID).First(&destAccount).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("destination account not found")
			}
			return err
		}
		
		// Lock accounts in consistent order to prevent deadlocks
		if sourceAccount.ID > destAccount.ID {
			// Re-lock in correct order
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&destAccount, destAccount.ID).Error; err != nil {
				return err
			}
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&sourceAccount, sourceAccount.ID).Error; err != nil {
				return err
			}
		}
		
		// Validate business rules using constraint functions
		if err := ValidateTransaction(sourceAccount.Balance, amount); err != nil {
			return err
		}
		
		// Update balances
		sourceAccount.Balance = sourceAccount.Balance.Sub(amount)
		destAccount.Balance = destAccount.Balance.Add(amount)
		
		// Save account updates
		if err := tx.Save(&sourceAccount).Error; err != nil {
			return err
		}
		
		if err := tx.Save(&destAccount).Error; err != nil {
			return err
		}
		
		// Create transaction record with database IDs
		transaction := &Transaction{
			SourceAccountID:      sourceAccount.ID,
			DestinationAccountID: destAccount.ID,
			Amount:               amount,
		}
		
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}
		
		return nil
	})
}

func (r *Repository) UpdateAccountBalance(accountID uint, balance decimal.Decimal) error {
	return r.db.Model(&Account{}).Where("account_id = ?", accountID).Update("balance", balance).Error
}