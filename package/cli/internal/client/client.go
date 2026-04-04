package client

import "time"

type Client interface {
	Login(email, password string) (*LoginResponse, error)
	Logout() error
	Whoami() (*User, error)

	GetUser(id string) (*User, error)
	ListUsers() ([]User, error)

	GetAccount(id string) (*Account, error)
	ListAccounts(userID string) ([]Account, error)

	GetTransaction(id string) (*Transaction, error)
	ListTransactions(userID string, limit, offset int) ([]Transaction, error)
	SimulateTransaction(params map[string]interface{}) (*Transaction, error)

	CreateTransfer(from, to string, amount int64, dryRun bool) (*Transfer, error)
	GetTransfer(id string) (*Transfer, error)

	GetLedgerAudit() (*LedgerAudit, error)
	GetLedgerEntries() ([]LedgerEntry, error)

	GetLogs(tail int) ([]LogEntry, error)
	DebugTransaction(id string) (*TransactionDebug, error)
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type Account struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AccountType string    `json:"account_type"`
	Balance     int64     `json:"balance"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Transfer struct {
	ID          string    `json:"id"`
	FromAccount string    `json:"from_account"`
	ToAccount   string    `json:"to_account"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type LedgerAudit struct {
	TotalEntries  int64     `json:"total_entries"`
	TotalDebits   int64     `json:"total_debits"`
	TotalCredits  int64     `json:"total_credits"`
	Balanced      bool      `json:"balanced"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type LedgerEntry struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Type      string    `json:"type"`
	Amount    int64     `json:"amount"`
	Balance   int64     `json:"balance"`
	Timestamp time.Time `json:"timestamp"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type TransactionDebug struct {
	Transaction *Transaction `json:"transaction"`
	Steps       []DebugStep  `json:"steps"`
}

type DebugStep struct {
	Step   string    `json:"step"`
	Time   time.Time `json:"time"`
	Status string    `json:"status"`
	Detail string    `json:"detail"`
}

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Login(email, password string) (*LoginResponse, error) {
	return &LoginResponse{
		AccessToken:  "mock_token_" + email,
		RefreshToken: "mock_refresh_" + email,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: &User{
			ID:        "user_001",
			Email:     email,
			Name:      "Test User",
			Active:    true,
			CreatedAt: time.Now(),
		},
	}, nil
}

func (m *MockClient) Logout() error {
	return nil
}

func (m *MockClient) Whoami() (*User, error) {
	return &User{
		ID:        "user_001",
		Email:     "admin@aetherbank.com",
		Name:      "Admin User",
		Active:    true,
		CreatedAt: time.Now(),
	}, nil
}

func (m *MockClient) GetUser(id string) (*User, error) {
	return &User{
		ID:        id,
		Email:     "user" + id + "@aetherbank.com",
		Name:      "User " + id,
		Active:    true,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}, nil
}

func (m *MockClient) ListUsers() ([]User, error) {
	return []User{
		{ID: "user_001", Email: "alice@aetherbank.com", Name: "Alice Smith", Active: true, CreatedAt: time.Now()},
		{ID: "user_002", Email: "bob@aetherbank.com", Name: "Bob Jones", Active: true, CreatedAt: time.Now()},
		{ID: "user_003", Email: "charlie@aetherbank.com", Name: "Charlie Brown", Active: false, CreatedAt: time.Now()},
	}, nil
}

func (m *MockClient) GetAccount(id string) (*Account, error) {
	return &Account{
		ID:          id,
		UserID:      "user_001",
		AccountType: "checking",
		Balance:     500000,
		Currency:    "EUR",
		Status:      "active",
		CreatedAt:   time.Now(),
	}, nil
}

func (m *MockClient) ListAccounts(userID string) ([]Account, error) {
	return []Account{
		{ID: "acc_001", UserID: userID, AccountType: "checking", Balance: 500000, Currency: "EUR", Status: "active", CreatedAt: time.Now()},
		{ID: "acc_002", UserID: userID, AccountType: "savings", Balance: 1000000, Currency: "EUR", Status: "active", CreatedAt: time.Now()},
		{ID: "acc_003", UserID: userID, AccountType: "investment", Balance: 2500000, Currency: "EUR", Status: "active", CreatedAt: time.Now()},
	}, nil
}

func (m *MockClient) GetTransaction(id string) (*Transaction, error) {
	return &Transaction{
		ID:          id,
		AccountID:   "acc_001",
		Type:        "credit",
		Status:      "completed",
		Amount:      10000,
		Currency:    "EUR",
		Description: "Salary deposit",
		CreatedAt:   time.Now().Add(-1 * time.Hour),
	}, nil
}

func (m *MockClient) ListTransactions(userID string, limit, offset int) ([]Transaction, error) {
	txs := []Transaction{
		{ID: "tx_001", AccountID: "acc_001", Type: "credit", Status: "completed", Amount: 50000, Currency: "EUR", Description: "Salary", CreatedAt: time.Now()},
		{ID: "tx_002", AccountID: "acc_001", Type: "debit", Status: "completed", Amount: -2500, Currency: "EUR", Description: "Rent", CreatedAt: time.Now().Add(-1 * time.Hour)},
		{ID: "tx_003", AccountID: "acc_001", Type: "card", Status: "completed", Amount: -1500, Currency: "EUR", Description: "Grocery", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{ID: "tx_004", AccountID: "acc_001", Type: "credit", Status: "completed", Amount: 10000, Currency: "EUR", Description: "Refund", CreatedAt: time.Now().Add(-3 * time.Hour)},
		{ID: "tx_005", AccountID: "acc_001", Type: "debit", Status: "pending", Amount: -5000, Currency: "EUR", Description: "Transfer", CreatedAt: time.Now().Add(-4 * time.Hour)},
	}
	if limit > 0 && limit < len(txs) {
		return txs[offset : offset+limit], nil
	}
	return txs, nil
}

func (m *MockClient) SimulateTransaction(params map[string]interface{}) (*Transaction, error) {
	return &Transaction{
		ID:          "tx_simulated",
		AccountID:   "acc_001",
		Type:        "credit",
		Status:      "simulated",
		Amount:      10000,
		Currency:    "EUR",
		Description: "Simulated transaction",
		CreatedAt:   time.Now(),
	}, nil
}

func (m *MockClient) CreateTransfer(from, to string, amount int64, dryRun bool) (*Transfer, error) {
	status := "completed"
	if dryRun {
		status = "dry_run"
	}
	return &Transfer{
		ID:          "transfer_001",
		FromAccount: from,
		ToAccount:   to,
		Amount:      amount,
		Status:      status,
		CreatedAt:   time.Now(),
	}, nil
}

func (m *MockClient) GetTransfer(id string) (*Transfer, error) {
	return &Transfer{
		ID:          id,
		FromAccount: "acc_001",
		ToAccount:   "acc_002",
		Amount:      10000,
		Status:      "completed",
		CreatedAt:   time.Now().Add(-30 * time.Minute),
	}, nil
}

func (m *MockClient) GetLedgerAudit() (*LedgerAudit, error) {
	return &LedgerAudit{
		TotalEntries:  10000,
		TotalDebits:   50000000,
		TotalCredits:  50000000,
		Balanced:      true,
		LastCheckedAt: time.Now(),
	}, nil
}

func (m *MockClient) GetLedgerEntries() ([]LedgerEntry, error) {
	return []LedgerEntry{
		{ID: "le_001", AccountID: "acc_001", Type: "debit", Amount: -5000, Balance: 495000, Timestamp: time.Now()},
		{ID: "le_002", AccountID: "acc_002", Type: "credit", Amount: 5000, Balance: 1005000, Timestamp: time.Now()},
		{ID: "le_003", AccountID: "acc_001", Type: "credit", Amount: 10000, Balance: 505000, Timestamp: time.Now().Add(-1 * time.Hour)},
	}, nil
}

func (m *MockClient) GetLogs(tail int) ([]LogEntry, error) {
	logs := []LogEntry{
		{Level: "INFO", Message: "Server started", Timestamp: time.Now().Add(-10 * time.Minute)},
		{Level: "INFO", Message: "User logged in", Timestamp: time.Now().Add(-5 * time.Minute)},
		{Level: "DEBUG", Message: "Transaction processed", Timestamp: time.Now().Add(-2 * time.Minute)},
		{Level: "WARN", Message: "Rate limit approaching", Timestamp: time.Now()},
		{Level: "ERROR", Message: "Connection timeout", Timestamp: time.Now()},
	}
	if tail > 0 && tail < len(logs) {
		return logs[len(logs)-tail:], nil
	}
	return logs, nil
}

func (m *MockClient) DebugTransaction(id string) (*TransactionDebug, error) {
	return &TransactionDebug{
		Transaction: &Transaction{
			ID:          id,
			AccountID:   "acc_001",
			Type:        "credit",
			Status:      "completed",
			Amount:      10000,
			Currency:    "EUR",
			Description: "Salary deposit",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
		Steps: []DebugStep{
			{Step: "created", Time: time.Now().Add(-1 * time.Hour), Status: "OK", Detail: "Transaction created"},
			{Step: "validated", Time: time.Now().Add(-55 * time.Minute), Status: "OK", Detail: "Validation passed"},
			{Step: "authorized", Time: time.Now().Add(-50 * time.Minute), Status: "OK", Detail: "Authorization granted"},
			{Step: "processed", Time: time.Now().Add(-45 * time.Minute), Status: "OK", Detail: "Processing complete"},
			{Step: "completed", Time: time.Now().Add(-40 * time.Minute), Status: "OK", Detail: "Transaction completed"},
		},
	}, nil
}
