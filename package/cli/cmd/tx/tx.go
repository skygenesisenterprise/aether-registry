package tx

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/output"

	"github.com/spf13/cobra"
)

var (
	txUserID string
	txLimit  int
	txOffset int
	jsonFlag bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions for a user",
	Long:  `Retrieve all transactions associated with a specific user.`,
	Example: `  bank tx list --user user_001
  bank tx list --user user_001 --json
  bank tx list --user user_001 --limit 10 --offset 0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		txs, err := apiClient.ListTransactions(txUserID, txLimit, txOffset)
		if err != nil {
			return fmt.Errorf("failed to list transactions: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(txs)
	},
}

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get transaction by ID",
	Long:  `Retrieve detailed information about a specific transaction.`,
	Example: `  bank tx get tx_001
  bank tx get tx_001 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		tx, err := apiClient.GetTransaction(args[0])
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(tx)
	},
}

var simulateCmd = &cobra.Command{
	Use:     "simulate",
	Short:   "Simulate a transaction",
	Long:    `Simulate a transaction with given parameters (dry run).`,
	Example: `  bank tx simulate --amount 1000 --from acc_001 --to acc_002`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		params := map[string]interface{}{
			"limit":  txLimit,
			"offset": txOffset,
		}
		tx, err := apiClient.SimulateTransaction(params)
		if err != nil {
			return fmt.Errorf("failed to simulate transaction: %w", err)
		}

		fmt.Println("Simulated transaction:")
		out := output.New(getOutputFormat())
		return out.Print(tx)
	},
}

var Command = &cobra.Command{
	Use:   "tx",
	Short: "Transaction commands",
	Long:  `Commands to manage and query transactions.`,
}

func getOutputFormat() string {
	if jsonFlag {
		return "json"
	}
	return "table"
}

func init() {
	listCmd.Flags().StringVarP(&txUserID, "user", "u", "", "User ID (required)")
	listCmd.Flags().IntVar(&txLimit, "limit", 0, "Limit number of results")
	listCmd.Flags().IntVar(&txOffset, "offset", 0, "Offset for pagination")
	listCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")
	listCmd.MarkFlagRequired("user")

	getCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	simulateCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	Command.AddCommand(listCmd)
	Command.AddCommand(getCmd)
	Command.AddCommand(simulateCmd)
}
