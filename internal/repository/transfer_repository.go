package repository

import (
	"context"
	"fmt"

	"github.com/dmsyudha/money-api/internal/domain"
	"github.com/dmsyudha/money-api/lib/bank"
	"gorm.io/gorm"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, fromAccountNumber string, toAccountNumber string, amount float64) (bool, error)
	HandleTransferCallback(transactionID string, status string) error
}

type transferRepository struct {
	db      *gorm.DB
	bankAPI *bank.BankAPI
}

func NewTransferRepository(db *gorm.DB, bankAPI *bank.BankAPI) TransferRepository {
	return &transferRepository{db: db, bankAPI: bankAPI}
}

func (r *transferRepository) CreateTransfer(ctx context.Context, fromAccountNumber string, toAccountNumber string, amount float64) (bool, error) {
	if valid, err := r.validateAccountsWithTimeout(ctx, fromAccountNumber, toAccountNumber); !valid {
		return false, err
	}

	if err := r.performTransfer(fromAccountNumber, toAccountNumber, amount); err != nil {
		return false, err
	}

	if err := r.storeTransferDetails(fromAccountNumber, toAccountNumber, amount); err != nil {
		return false, err
	}

	return true, nil
}

func (r *transferRepository) HandleTransferCallback(transactionID string, status string) error {
	go func() {
		err := r.bankAPI.Callback(transactionID, status)
		if err != nil {
			fmt.Printf("Error handling transfer callback: %v\n", err)
		}
	}()
	return nil
}

func (r *transferRepository) validateAccountsWithTimeout(ctx context.Context, fromAccountNumber, toAccountNumber string) (bool, error) {
	results := make(chan bool, 2)
	errors := make(chan error, 2)

	validate := func(accountNumber string) {
		select {
		case <-ctx.Done():
			errors <- ctx.Err()
		default:
			valid, err := r.bankAPI.Validate(accountNumber)
			if err != nil {
				errors <- fmt.Errorf("error validating account %s: %w", accountNumber, err)
				return
			}
			results <- valid
		}
	}

	go validate(fromAccountNumber)
	go validate(toAccountNumber)

	validCount := 0
	for i := 0; i < 2; i++ {
		select {
		case err := <-errors:
			return false, err
		case valid := <-results:
			if valid {
				validCount++
			}
		}
	}

	return validCount == 2, nil
}

func (r *transferRepository) performTransfer(fromAccountNumber, toAccountNumber string, amount float64) error {
	return r.bankAPI.Transfer(fromAccountNumber, toAccountNumber, amount)
}

func (r *transferRepository) storeTransferDetails(fromAccountNumber, toAccountNumber string, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var fromAccount, toAccount domain.Account

		if err := tx.Where("account_number = ?", fromAccountNumber).First(&fromAccount).Error; err != nil {
			return err
		}

		if err := tx.Where("account_number = ?", toAccountNumber).First(&toAccount).Error; err != nil {
			return err
		}

		transfer := domain.Transfer{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        amount,
		}

		return tx.Create(&transfer).Error
	})
}
