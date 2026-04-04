package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0"
var buildDate = "unknown"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print version information",
	Long:    `Display version and build information for the bank CLI.`,
	Example: `  bank version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("bank version %s\n", version)
		fmt.Printf("Build date: %s\n", buildDate)
	},
}

func PrintVersion() {
	fmt.Printf("bank version %s\n", version)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
