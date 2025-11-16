package models

import (
	"time"
)

type Account struct {
	AccountID int64     `gorm:"primaryKey"`
	Balance   string    `gorm:"type:numeric(20,5);not null;check:balance >= 0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type CreateAccountRequest struct {
	AccountID      int64  `json:"account_id" validate:"required,gt=0"`
	InitialBalance string `json:"initial_balance" validate:"required"`
}

type AccountResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}
