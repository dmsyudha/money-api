package handler

import (
	"fmt"
	"net/http"

	"github.com/dmsyudha/money-api/internal/service"
)

type AccountHandler interface {
	ValidateAccountHandler() http.HandlerFunc
}

type accountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

func (h *accountHandler) ValidateAccountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateRequestMethod(w, r) {
			return
		}

		accountNumber, accountName, valid := validateRequestParams(w, r)
		if !valid {
			return
		}

		if !processValidation(w, accountNumber, accountName, h) {
			return
		}

		successResponse(w)
	}
}

func validateRequestMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func validateRequestParams(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	accountNumber := r.URL.Query().Get("accountNumber")
	accountName := r.URL.Query().Get("accountName")

	if accountNumber == "" || accountName == "" {
		http.Error(w, "Missing account number or account name", http.StatusBadRequest)
		return "", "", false
	}
	return accountNumber, accountName, true
}

func processValidation(w http.ResponseWriter, accountNumber string, accountName string, h *accountHandler) bool {
	valid, err := h.accountService.ValidateAccount(accountNumber, accountName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validating account: %v", err), http.StatusInternalServerError)
		return false
	}

	if !valid {
		http.Error(w, "Invalid account details", http.StatusUnauthorized)
		return false
	}
	return true
}

func successResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Account validated successfully")
}