package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewHTTPClient(baseURL, token string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *HTTPClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *HTTPClient) Login(email, password string) (*LoginResponse, error) {
	respBody, err := c.doRequest("POST", "/api/v1/auth/login", loginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	var authResp struct {
		Success bool           `json:"success"`
		Data    *LoginResponse `json:"data"`
		Error   string         `json:"error"`
	}

	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !authResp.Success {
		return nil, fmt.Errorf("login failed: %s", authResp.Error)
	}

	return authResp.Data, nil
}

func (c *HTTPClient) Logout() error {
	_, err := c.doRequest("POST", "/api/v1/auth/logout", nil)
	return err
}

func (c *HTTPClient) Whoami() (*User, error) {
	respBody, err := c.doRequest("GET", "/api/v1/auth/account/me", nil)
	if err != nil {
		return nil, err
	}

	var authResp struct {
		Success bool `json:"success"`
		Data    struct {
			User *User `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return authResp.Data.User, nil
}

func (c *HTTPClient) GetUser(id string) (*User, error) {
	respBody, err := c.doRequest("GET", "/api/v1/admin/users/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			User *User `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data.User, nil
}

func (c *HTTPClient) ListUsers() ([]User, error) {
	respBody, err := c.doRequest("GET", "/api/v1/admin/users", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool   `json:"success"`
		Data    []User `json:"data"`
		Total   int    `json:"total"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) GetAccount(id string) (*Account, error) {
	respBody, err := c.doRequest("GET", "/api/v1/accounts/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool     `json:"success"`
		Data    *Account `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) ListAccounts(userID string) ([]Account, error) {
	respBody, err := c.doRequest("GET", "/api/v1/accounts", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool      `json:"success"`
		Data    []Account `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) GetTransaction(id string) (*Transaction, error) {
	respBody, err := c.doRequest("GET", "/api/v1/transactions/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    *Transaction `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) ListTransactions(userID string, limit, offset int) ([]Transaction, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}

	path := "/api/v1/transactions"
	if query := params.Encode(); query != "" {
		path += "?" + query
	}

	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool          `json:"success"`
		Data    []Transaction `json:"data"`
		Total   int           `json:"total"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) SimulateTransaction(params map[string]interface{}) (*Transaction, error) {
	respBody, err := c.doRequest("POST", "/api/v1/transactions/simulate", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    *Transaction `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

type transferRequest struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency,omitempty"`
}

func (c *HTTPClient) CreateTransfer(from, to string, amount int64, dryRun bool) (*Transfer, error) {
	req := transferRequest{
		FromAccount: from,
		ToAccount:   to,
		Amount:      amount,
		Currency:    "EUR",
	}

	path := "/api/v1/transfers"
	if dryRun {
		path += "/simulate"
	}

	respBody, err := c.doRequest("POST", path, req)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool      `json:"success"`
		Data    *Transfer `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) GetTransfer(id string) (*Transfer, error) {
	respBody, err := c.doRequest("GET", "/api/v1/transfers/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool      `json:"success"`
		Data    *Transfer `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) GetLedgerAudit() (*LedgerAudit, error) {
	respBody, err := c.doRequest("GET", "/api/v1/ledger/balance", nil)
	if err != nil {
		return nil, err
	}

	var ledger struct {
		Success      bool  `json:"success"`
		TotalDebits  int64 `json:"total_debits"`
		TotalCredits int64 `json:"total_credits"`
	}

	if err := json.Unmarshal(respBody, &ledger); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &LedgerAudit{
		TotalDebits:   ledger.TotalDebits,
		TotalCredits:  ledger.TotalCredits,
		Balanced:      ledger.TotalDebits == ledger.TotalCredits,
		LastCheckedAt: time.Now(),
	}, nil
}

func (c *HTTPClient) GetLedgerEntries() ([]LedgerEntry, error) {
	respBody, err := c.doRequest("GET", "/api/v1/ledger/entries", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool          `json:"success"`
		Data    []LedgerEntry `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *HTTPClient) GetLogs(tail int) ([]LogEntry, error) {
	respBody, err := c.doRequest("GET", "/api/v1/system-logs", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Success bool       `json:"success"`
		Data    []LogEntry `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	logs := resp.Data
	if tail > 0 && tail < len(logs) {
		logs = logs[len(logs)-tail:]
	}

	return logs, nil
}

func (c *HTTPClient) DebugTransaction(id string) (*TransactionDebug, error) {
	tx, err := c.GetTransaction(id)
	if err != nil {
		return nil, err
	}

	return &TransactionDebug{
		Transaction: tx,
		Steps: []DebugStep{
			{Step: "fetched", Time: time.Now(), Status: "OK", Detail: "Transaction fetched from API"},
		},
	}, nil
}

func NewClient(baseURL, token string) Client {
	return NewHTTPClient(baseURL, token)
}
