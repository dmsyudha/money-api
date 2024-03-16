package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	api "github.com/dmsyudha/money-api/pkg/http_client"
)

var (
	baseURL = os.Getenv("BANK_BASE_URL")
	token   = os.Getenv("BANK_TOKEN")
)

type BankAPI struct {
	client *api.APIClient
}

func NewBankAPI(client *api.APIClient) *BankAPI {
	return &BankAPI{
		client: client,
	}
}

func (b *BankAPI) Validate(accountNumber string) (bool, error) {
	url := fmt.Sprintf("%s/validate", baseURL)
	data := map[string]interface{}{
		"account_number": accountNumber,
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	req, err := api.NewRequest("POST", url, nil, headers, data)
	if err != nil {
		return false, err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to validate account")
	}

	var result struct {
		Valid bool `json:"valid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Valid, nil
}

func (b *BankAPI) Transfer(fromAccount, toAccount string, amount float64) error {
	url := fmt.Sprintf("%s/transfer", baseURL)
	data := map[string]interface{}{
		"from_account": fromAccount,
		"to_account":   toAccount,
		"amount":       amount,
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	req, err := api.NewRequest("POST", url, nil, headers, data)
	if err != nil {
		return err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to transfer money")
	}

	return nil
}

func (b *BankAPI) Callback(transactionID string, status string) error {
	url := fmt.Sprintf("%s/callback", baseURL)
	data := map[string]interface{}{
		"transaction_id": transactionID,
		"status":         status,
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	req, err := api.NewRequest("POST", url, nil, headers, data)
	if err != nil {
		return err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send callback")
	}

	return nil
}
