package ledger

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/output"

	"github.com/spf13/cobra"
)

var jsonFlag bool

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run ledger audit",
	Long:  `Verify ledger integrity by checking that total debits equal total credits.`,
	Example: `  bank ledger audit
  bank ledger audit --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		audit, err := apiClient.GetLedgerAudit()
		if err != nil {
			return fmt.Errorf("failed to run ledger audit: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(audit)
	},
}

var entriesCmd = &cobra.Command{
	Use:   "entries",
	Short: "List ledger entries",
	Long:  `Retrieve recent ledger entries.`,
	Example: `  bank ledger entries
  bank ledger entries --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		entries, err := apiClient.GetLedgerEntries()
		if err != nil {
			return fmt.Errorf("failed to get ledger entries: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(entries)
	},
}

var Command = &cobra.Command{
	Use:   "ledger",
	Short: "Ledger commands",
	Long:  `Commands to query and audit the ledger.`,
}

func getOutputFormat() string {
	if jsonFlag {
		return "json"
	}
	return "table"
}

func init() {
	auditCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")
	entriesCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	Command.AddCommand(auditCmd)
	Command.AddCommand(entriesCmd)
}
