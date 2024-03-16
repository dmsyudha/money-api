package repository

import (
    "github.com/dmsyudha/money-api/internal/domain"
    "github.com/dmsyudha/money-api/lib/bank"
    "gorm.io/gorm"
)

type AccountRepository interface {
    ValidateAccount(accountNumber string, accountName string) (bool, error)
}

type accountRepository struct {
    db *gorm.DB
    bankAPI *bank.BankAPI
}

func NewAccountRepository(db *gorm.DB, bankAPI *bank.BankAPI) AccountRepository {
    return &accountRepository{db: db, bankAPI: bankAPI}
}

func (r *accountRepository) ValidateAccount(accountNumber string, accountName string) (bool, error) {
    isValid, err := r.bankAPI.Validate(accountNumber)
    if err != nil {
        return false, err
    }
    if !isValid {
        return false, nil
    }

    var account domain.Account
    if result := r.db.Where("account_number = ? AND account_name = ?", accountNumber, accountName).First(&account); result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return false, nil
        }
        return false, result.Error
    }

    return true, nil
}
