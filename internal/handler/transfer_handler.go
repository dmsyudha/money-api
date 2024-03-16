package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmsyudha/money-api/internal/domain"
	"github.com/dmsyudha/money-api/internal/service"
)

type TransferHandler interface {
	CreateTransferHandler() http.HandlerFunc
}

type transferHandler struct {
	transferService service.TransferService
}

func NewTransferHandler(transferService service.TransferService) TransferHandler {
	return &transferHandler{
		transferService: transferService,
	}
}

func (h *transferHandler) CreateTransferHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var transferRequest domain.Transfer
		err := json.NewDecoder(r.Body).Decode(&transferRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		transferResponse, err := h.transferService.CreateTransfer(transferRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating transfer: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transferResponse)
	}
}
func (h *transferHandler) HandleTransferCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var callbackRequest struct {
			TransactionID string `json:"transactionId"`
			Status        string `json:"status"`
		}
		err := json.NewDecoder(r.Body).Decode(&callbackRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = h.transferService.HandleTransferCallback(callbackRequest.TransactionID, callbackRequest.Status)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error handling transfer callback: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Callback handled successfully"})
	}
}
