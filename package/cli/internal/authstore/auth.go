package authstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skygenesisenterprise/aether-bank/cli/internal/client"
	"github.com/spf13/viper"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user,omitempty"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func SaveToken(token *Token) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configDir := filepath.Join(home, ".bank")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	tokenFile := filepath.Join(configDir, "token.json")
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal token: %w", err)
	}

	if err := os.WriteFile(tokenFile, data, 0600); err != nil {
		return fmt.Errorf("could not write token file: %w", err)
	}

	viper.Set("token", token.AccessToken)
	configFile := filepath.Join(configDir, "config.yaml")
	viper.SetConfigFile(configFile)
	viper.Set("token", token.AccessToken)
	_ = viper.WriteConfig()

	return nil
}

func LoadToken() (*Token, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	tokenFile := filepath.Join(home, ".bank", "token.json")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("could not read token file: %w", err)
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("could not unmarshal token: %w", err)
	}

	return &token, nil
}

func ClearToken() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	tokenFile := filepath.Join(home, ".bank", "token.json")
	if err := os.Remove(tokenFile); err != nil {
		return fmt.Errorf("could not remove token file: %w", err)
	}

	viper.Set("token", "")
	configFile := filepath.Join(home, ".bank", "config.yaml")
	viper.SetConfigFile(configFile)
	viper.Set("token", "")
	_ = viper.WriteConfig()

	return nil
}

func GetToken() (string, error) {
	token, err := LoadToken()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func UserFromClient(u *client.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
