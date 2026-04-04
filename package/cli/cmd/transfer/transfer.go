package transfer

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/output"

	"github.com/spf13/cobra"
)

var (
	fromAccount string
	toAccount   string
	amount      int64
	dryRun      bool
	jsonFlag    bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a transfer",
	Long:  `Create a new transfer between accounts. Use --dry-run to simulate without executing.`,
	Example: `  bank transfer create --from acc_001 --to acc_002 --amount 10000
  bank transfer create --from acc_001 --to acc_002 --amount 10000 --dry-run
  bank transfer create --from acc_001 --to acc_002 --amount 10000 --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		transfer, err := apiClient.CreateTransfer(fromAccount, toAccount, amount, dryRun)
		if err != nil {
			return fmt.Errorf("failed to create transfer: %w", err)
		}

		if dryRun {
			fmt.Println("Dry run - transfer not executed")
		}
		out := output.New(getOutputFormat())
		return out.Print(transfer)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status <id>",
	Short: "Get transfer status",
	Long:  `Retrieve the status of a specific transfer.`,
	Example: `  bank transfer status transfer_001
  bank transfer status transfer_001 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		transfer, err := apiClient.GetTransfer(args[0])
		if err != nil {
			return fmt.Errorf("failed to get transfer: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(transfer)
	},
}

var Command = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer commands",
	Long:  `Commands to manage and execute transfers between accounts.`,
}

func getOutputFormat() string {
	if jsonFlag {
		return "json"
	}
	return "table"
}

func init() {
	createCmd.Flags().StringVarP(&fromAccount, "from", "f", "", "Source account ID (required)")
	createCmd.Flags().StringVarP(&toAccount, "to", "t", "", "Destination account ID (required)")
	createCmd.Flags().Int64VarP(&amount, "amount", "a", 0, "Amount in cents (required)")
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate without executing")
	createCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")
	createCmd.MarkFlagRequired("from")
	createCmd.MarkFlagRequired("to")
	createCmd.MarkFlagRequired("amount")

	statusCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	Command.AddCommand(createCmd)
	Command.AddCommand(statusCmd)
}
