package repository

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	ID        uint            `gorm:"primarykey" json:"-"`
	AccountID uint            `gorm:"uniqueIndex;not null" json:"account_id"`
	Balance   decimal.Decimal `gorm:"type:decimal(20,5);not null;check:balance >= 0" json:"balance"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
}

type Transaction struct {
	ID                   uint            `gorm:"primarykey" json:"id"`
	SourceAccountID      uint            `gorm:"not null;index" json:"source_account_id"`
	DestinationAccountID uint            `gorm:"not null;index" json:"destination_account_id"`
	Amount               decimal.Decimal `gorm:"type:decimal(20,5);not null;check:amount > 0" json:"amount"`
	CreatedAt            time.Time       `gorm:"index" json:"created_at"`
	UpdatedAt            time.Time       `json:"-"`
	DeletedAt            gorm.DeletedAt  `gorm:"index" json:"-"`

	SourceAccount      Account `gorm:"foreignKey:SourceAccountID;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
	DestinationAccount Account `gorm:"foreignKey:DestinationAccountID;references:ID;constraint:OnDelete:RESTRICT" json:"-"`
}

func (Account) TableName() string {
	return "accounts"
}

func (Transaction) TableName() string {
	return "transactions"
}