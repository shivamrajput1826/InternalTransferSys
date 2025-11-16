package models

import "time"

type CreateTransactionRequest struct {
	SourceAccountID      int64  `json:"source_account_id" validate:"required,gt=0"`
	DestinationAccountID int64  `json:"destination_account_id" validate:"required,gt=0,nefield=SourceAccountID"`
	Amount               string `json:"amount" validate:"required"`
}

type Transaction struct {
	ID                   int64     `gorm:"primaryKey;autoIncrement"`
	SourceAccountID      int64     `gorm:"not null;index"`
	DestinationAccountID int64     `gorm:"not null;index"`
	Amount               string    `gorm:"type:numeric(20,5);not null;check:amount > 0"`
	Status               string    `gorm:"type:varchar(20);not null;default:'completed'"`
	CreatedAt            time.Time `gorm:"autoCreateTime;index"`
}
