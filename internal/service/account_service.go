package service

import (
	"github.com/dmsyudha/money-api/internal/repository"
)

type AccountService interface {
	ValidateAccount(accountNumber string, accountName string) (bool, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) ValidateAccount(accountNumber string, accountName string) (bool, error) {
	return s.accountRepo.ValidateAccount(accountNumber, accountName)
}
