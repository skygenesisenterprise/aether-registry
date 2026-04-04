package user

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/output"

	"github.com/spf13/cobra"
)

var jsonFlag bool

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get user by ID",
	Long:  `Retrieve detailed information about a specific user.`,
	Example: `  bank user get user_001
  bank user get user_001 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		user, err := apiClient.GetUser(args[0])
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(user)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `Retrieve a list of all users in the system.`,
	Example: `  bank user list
  bank user list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := client.NewMockClient()
		users, err := apiClient.ListUsers()
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		out := output.New(getOutputFormat())
		return out.Print(users)
	},
}

var Command = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  `Commands to manage and query user information.`,
}

func getOutputFormat() string {
	if jsonFlag {
		return "json"
	}
	return "table"
}

func init() {
	getCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")
	listCmd.Flags().BoolVar(&jsonFlag, "json", false, "JSON output")

	Command.AddCommand(getCmd)
	Command.AddCommand(listCmd)
}
