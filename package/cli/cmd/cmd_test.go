package cmd

import (
	"testing"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
)

type MockAPIClient struct {
	Users    []client.User
	Accounts []client.Account
}

func TestUserListCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	users, err := apiClient.ListUsers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

func TestUserGetCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	user, err := apiClient.GetUser("user_001")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.ID != "user_001" {
		t.Errorf("Expected user_001, got %s", user.ID)
	}
}

func TestAccountListCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	accounts, err := apiClient.ListAccounts("user_001")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(accounts) != 3 {
		t.Errorf("Expected 3 accounts, got %d", len(accounts))
	}
}

func TestTransactionListCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	txs, err := apiClient.ListTransactions("user_001", 0, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(txs) != 5 {
		t.Errorf("Expected 5 transactions, got %d", len(txs))
	}
}

func TestTransferCreateCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	transfer, err := apiClient.CreateTransfer("acc_001", "acc_002", 10000, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if transfer.Amount != 10000 {
		t.Errorf("Expected amount 10000, got %d", transfer.Amount)
	}
	if transfer.Status != "completed" {
		t.Errorf("Expected status completed, got %s", transfer.Status)
	}
}

func TestTransferDryRunCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	transfer, err := apiClient.CreateTransfer("acc_001", "acc_002", 10000, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if transfer.Status != "dry_run" {
		t.Errorf("Expected status dry_run, got %s", transfer.Status)
	}
}

func TestLedgerAuditCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	audit, err := apiClient.GetLedgerAudit()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !audit.Balanced {
		t.Error("Expected ledger to be balanced")
	}
}

func TestDebugTransactionCommand(t *testing.T) {
	apiClient := client.NewMockClient()
	debug, err := apiClient.DebugTransaction("tx_001")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(debug.Steps) != 5 {
		t.Errorf("Expected 5 debug steps, got %d", len(debug.Steps))
	}
}
