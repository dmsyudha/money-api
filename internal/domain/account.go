package domain

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Account represents a bank account in the system.
type Account struct {
    gorm.Model
    ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
    AccountNumber string    `gorm:"type:varchar(100);uniqueIndex"`
    AccountName   string    `gorm:"type:varchar(100)"`
    BankName      string    `gorm:"type:varchar(100)"`
    Balance       float64   `gorm:"type:numeric"`
}
