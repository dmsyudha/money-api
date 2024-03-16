package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	httpClient *http.Client
}

func NewAPIClient(httpClient *http.Client) *APIClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &APIClient{httpClient: httpClient}
}

func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func NewRequest(method, url string, params map[string]string, headers map[string]string, data interface{}) (*http.Request, error) {
	var body []byte
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = jsonData
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
