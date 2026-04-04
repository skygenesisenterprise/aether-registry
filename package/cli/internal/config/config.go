package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Env          string                       `mapstructure:"env"`
	APIURL       string                       `mapstructure:"api_url"`
	Token        string                       `mapstructure:"token"`
	Output       string                       `mapstructure:"output"`
	Debug        bool                         `mapstructure:"debug"`
	UseMock      bool                         `mapstructure:"use_mock"`
	Timeout      int                          `mapstructure:"timeout"`
	Environments map[string]EnvironmentConfig `mapstructure:"environments"`
}

type EnvironmentConfig struct {
	APIURL string `mapstructure:"api_url"`
}

func Default() *Config {
	return &Config{
		Env:     "staging",
		APIURL:  "http://localhost:8080",
		Output:  "table",
		Debug:   false,
		UseMock: true,
		Timeout: 30,
		Environments: map[string]EnvironmentConfig{
			"staging": {APIURL: "http://localhost:8080"},
			"prod":    {APIURL: "https://api.bank.skygenesisenterprise.com"},
			"local":   {APIURL: "http://localhost:8080"},
		},
	}
}

func New(env string) (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(home, ".bank", "config.yaml")
	viper.SetConfigFile(configPath)
	viper.SetDefault("env", env)
	viper.SetDefault("api_url", "http://localhost:8080")
	viper.SetDefault("output", "table")
	viper.SetDefault("debug", false)
	viper.SetDefault("use_mock", true)
	viper.SetDefault("timeout", 30)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		Env:     viper.GetString("env"),
		APIURL:  viper.GetString("api_url"),
		Token:   viper.GetString("token"),
		Output:  viper.GetString("output"),
		Debug:   viper.GetBool("debug"),
		UseMock: viper.GetBool("use_mock"),
		Timeout: viper.GetInt("timeout"),
	}

	if envConfigs, ok := viper.Get("environments").(map[string]interface{}); ok {
		cfg.Environments = make(map[string]EnvironmentConfig)
		for k, v := range envConfigs {
			if m, ok := v.(map[string]interface{}); ok {
				apiURL, _ := m["api_url"].(string)
				cfg.Environments[k] = EnvironmentConfig{APIURL: apiURL}
			}
		}
	}

	if cfg.Environments == nil {
		cfg.Environments = Default().Environments
	}

	if envCfg, ok := cfg.Environments[env]; ok && envCfg.APIURL != "" {
		cfg.APIURL = envCfg.APIURL
	}

	return cfg, nil
}

func LoadConfig(configFile string) error {
	return nil
}

func (c *Config) SaveToken(token string) error {
	c.Token = token

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configDir := filepath.Join(home, ".bank")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	viper.SetConfigFile(filepath.Join(configDir, "config.yaml"))
	viper.Set("token", token)

	return viper.WriteConfig()
}

func (c *Config) ClearToken() error {
	c.Token = ""

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	viper.SetConfigFile(filepath.Join(home, ".bank", "config.yaml"))
	viper.Set("token", "")

	return viper.WriteConfig()
}

func (c *Config) GetAPIURL() string {
	return c.APIURL
}

func (c *Config) GetToken() string {
	return c.Token
}

func (c *Config) IsMockEnabled() bool {
	return c.UseMock
}
