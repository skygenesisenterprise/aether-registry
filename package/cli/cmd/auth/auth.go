package auth

import (
	"fmt"
	"os"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/authstore"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/skygenesisenterprise/aether-bank/cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	email    string
	password string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the banking API",
	Long:  `Authenticate with the banking API and store the token for future use.`,
	Example: `  bank auth login --email admin@aetherbank.com --password secret
  bank auth login --env prod --email user@bank.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()

		var resp *client.LoginResponse
		var err error

		if cfg.IsMockEnabled() {
			mockClient := client.NewMockClient()
			resp, err = mockClient.Login(email, password)
		} else {
			httpClient := client.NewHTTPClient(cfg.GetAPIURL(), "")
			resp, err = httpClient.Login(email, password)
		}

		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if err := authstore.SaveToken(&authstore.Token{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			TokenType:    resp.TokenType,
			ExpiresIn:    resp.ExpiresIn,
			User:         authstore.UserFromClient(resp.User),
		}); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Printf("Logged in as %s (%s)\n", resp.User.Name, resp.User.Email)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "Logout and clear stored credentials",
	Long:    `Clear the stored authentication token.`,
	Example: `  bank auth logout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()

		if !cfg.IsMockEnabled() {
			token, _ := authstore.LoadToken()
			if token != nil {
				httpClient := client.NewHTTPClient(cfg.GetAPIURL(), token.AccessToken)
				_ = httpClient.Logout()
			}
		}

		if err := authstore.ClearToken(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}
		fmt.Println("Logged out successfully")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Short:   "Show current authenticated user",
	Long:    `Display information about the currently authenticated user.`,
	Example: `  bank auth whoami`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := authstore.LoadToken()
		if err != nil {
			return fmt.Errorf("not authenticated: %w", err)
		}

		if token.User != nil {
			fmt.Printf("User: %s\n", token.User.Name)
			fmt.Printf("Email: %s\n", token.User.Email)
			fmt.Printf("ID: %s\n", token.User.ID)
			return nil
		}

		cfg := getConfig()

		var user *client.User
		if cfg.IsMockEnabled() {
			mockClient := client.NewMockClient()
			user, err = mockClient.Whoami()
		} else {
			httpClient := client.NewHTTPClient(cfg.GetAPIURL(), token.AccessToken)
			user, err = httpClient.Whoami()
		}

		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		fmt.Printf("User: %s\n", user.Name)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("ID: %s\n", user.ID)
		return nil
	},
}

AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage authentication and credentials for the banking CLI.`,
}

func getConfig() *config.Config {
	cfg, err := config.New("staging")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config: %v\n", err)
		return config.Default()
	}
	return cfg
}

func init() {
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email address (required)")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password (required)")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(whoamiCmd)
}
