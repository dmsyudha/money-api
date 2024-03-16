package service

import (
	"context"
	"errors"

	"github.com/dmsyudha/money-api/internal/domain"
	"github.com/dmsyudha/money-api/internal/repository"
)

type TransferService interface {
	CreateTransfer(transfer domain.Transfer) (domain.Transfer, error)
	HandleTransferCallback(transactionID string, status string) error
}

type transferService struct {
	transferRepo repository.TransferRepository
	accountRepo  repository.AccountRepository
}

func NewTransferService(transferRepo repository.TransferRepository, accountRepo repository.AccountRepository) TransferService {
	return &transferService{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
	}
}

func (s *transferService) CreateTransfer(transfer domain.Transfer) (domain.Transfer, error) {
	// Validate both accounts before initiating a transfer
	validFromAccount, err := s.accountRepo.ValidateAccount(transfer.FromAccountID.String(), "")
	if err != nil {
		return domain.Transfer{}, err
	}
	if !validFromAccount {
		return domain.Transfer{}, errors.New("from account is invalid")
	}

	validToAccount, err := s.accountRepo.ValidateAccount(transfer.ToAccountID.String(), "")
	if err != nil {
		return domain.Transfer{}, err
	}
	if !validToAccount {
		return domain.Transfer{}, errors.New("to account is invalid")
	}

	// Perform the transfer through the repository
	success, err := s.transferRepo.CreateTransfer(context.Background(), transfer.FromAccountID.String(), transfer.ToAccountID.String(), transfer.Amount)
	if err != nil {
		return domain.Transfer{}, err
	}
	if !success {
		return domain.Transfer{}, errors.New("transfer failed")
	}

	return transfer, nil
}

func (s *transferService) HandleTransferCallback(transactionID string, status string) error {
	return s.transferRepo.HandleTransferCallback(transactionID, status)
}
