package debug

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/output"

	"github.com/spf13/cobra"
)

var (
	tailLines int
	jsonFlag  bool
)

var logsCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail recent logs",
	Long:  `Display recent log entries from the banking system.`,
	Example: `  bank logs tail
  bank logs tail --lines 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		logs, err := apiClient.GetLogs(tailLines)
		if err != nil {
			return fmt.Errorf("failed to get logs: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(logs)
	},
}

var debugTxCmd = &cobra.Command{
	Use:   "tx <id>",
	Short: "Debug a transaction",
	Long:  `Get detailed debug information about a transaction including all processing steps.`,
	Example: `  bank debug tx tx_001
  bank debug tx tx_001 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		debug, err := apiClient.DebugTransaction(args[0])
		if err != nil {
			return fmt.Errorf("failed to debug transaction: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(debug)
	},
}

var logsSubCmd = &cobra.Command{
	Use:   "logs",
	Short: "Log commands",
	Long:  `Commands to view system logs.`,
}

var Command = &cobra.Command{
	Use:   "debug",
	Short: "Debug commands",
	Long:  `Commands for debugging and diagnostics.`,
}

func getOutputFormat() string {
	if jsonFlag {
		return "json"
	}
	return "table"
}

func init() {
	logsCmd.Flags().IntVar(&tailLines, "lines", 10, "Number of lines to show")
	logsCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	debugTxCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	logsSubCmd.AddCommand(logsCmd)
	Command.AddCommand(logsSubCmd)
	Command.AddCommand(debugTxCmd)
}
