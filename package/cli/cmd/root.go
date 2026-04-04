package cmd

import (
	"fmt"
	"os"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "bank",
	Short: "Aether Bank CLI - Internal banking operations tool",
	Long: `Aether Bank CLI is an internal tool for developers and employees 
to interact with the banking backend in a stable, secure, and scriptable way.

Environment:
  --config, -c    Config file path (default ~/.bank/config.yaml)
  --env, -e      Environment: local, staging, prod

Examples:
  bank auth login
  bank user list
  bank account list --user <id>
  bank tx list --user <id> --json`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.LoadConfig(cfgFile)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default ~/.bank/config.yaml)")
	rootCmd.PersistentFlags().StringP("env", "e", "local", "environment: local, staging, prod")
	rootCmd.PersistentFlags().Bool("json", false, "JSON output")
	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(txCmd)
	rootCmd.AddCommand(transferCmd)
	rootCmd.AddCommand(ledgerCmd)
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	var err error
	env, _ := rootCmd.Flags().GetString("env")
	cfg, err = config.New(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config: %v\n", err)
		cfg = config.Default()
	}
}
