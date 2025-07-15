package repository

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// BeforeUpdate hook for Account
func (a *Account) BeforeUpdate(tx *gorm.DB) error {
	if a.Balance.IsNegative() {
		return errors.New("balance cannot be negative")
	}
	return nil
}

// BeforeCreate hook for Transaction
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("transaction amount must be positive")
	}
	
	if t.SourceAccountID == t.DestinationAccountID {
		return errors.New("source and destination accounts cannot be the same")
	}
	
	if t.SourceAccountID == 0 || t.DestinationAccountID == 0 {
		return errors.New("account IDs must be valid")
	}
	
	return nil
}

// BeforeCreate hook for Account
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.AccountID == 0 {
		return errors.New("account_id must be a positive integer")
	}
	
	if a.Balance.IsNegative() {
		return errors.New("balance cannot be negative")
	}
	
	return nil
}

// ValidateTransaction validates transaction business rules
func ValidateTransaction(sourceBalance, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}
	
	if sourceBalance.LessThan(amount) {
		return errors.New("insufficient balance")
	}
	
	return nil
}