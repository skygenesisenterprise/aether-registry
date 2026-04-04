package account

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/output"

	"github.com/spf13/cobra"
)

var (
	accountID string
	userID    string
	jsonFlag  bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts for a user",
	Long:  `Retrieve all accounts associated with a specific user.`,
	Example: `  bank account list --user user_001
  bank account list --user user_001 --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		accounts, err := apiClient.ListAccounts(userID)
		if err != nil {
			return fmt.Errorf("failed to list accounts: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(accounts)
	},
}

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get account by ID",
	Long:  `Retrieve detailed information about a specific account.`,
	Example: `  bank account get acc_001
  bank account get acc_001 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		account, err := apiClient.GetAccount(args[0])
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(account)
	},
}

var Command = &cobra.Command{
	Use:   "account",
	Short: "Account management commands",
	Long:  `Commands to manage and query account information.`,
}

func getOutputFormat() string {
	if jsonFlag {
		return "json"
	}
	return "table"
}

func init() {
	listCmd.Flags().StringVarP(&userID, "user", "u", "", "User ID (required)")
	listCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")
	listCmd.MarkFlagRequired("user")

	getCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	Command.AddCommand(listCmd)
	Command.AddCommand(getCmd)
}
