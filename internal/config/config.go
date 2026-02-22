package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DatabasePath string `mapstructure:"database_path"`
	MaxSuggestions int `mapstructure:"max_suggestions"`
	RecencyWeight float64 `mapstructure:"recency_weight"`
	FrequencyWeight float64 `mapstructure:"frequency_weight"`
	DirectoryWeight float64 `mapstructure:"directory_weight"`
	EnableColors bool `mapstructure:"enable_colors"`
	ShellIntegration string `mapstructure:"shell_integration"`
}

var cfg *Config

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".kwik-cmd")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("database_path", filepath.Join(homeDir, ".kwik-cmd", "commands.db"))
	viper.SetDefault("max_suggestions", 10)
	viper.SetDefault("recency_weight", 0.4)
	viper.SetDefault("frequency_weight", 0.4)
	viper.SetDefault("directory_weight", 0.2)
	viper.SetDefault("enable_colors", true)
	viper.SetDefault("shell_integration", "auto")

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		// Config doesn't exist, create default
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg = &Config{
				DatabasePath:    filepath.Join(homeDir, ".kwik-cmd", "commands.db"),
				MaxSuggestions: 10,
				RecencyWeight:   0.4,
				FrequencyWeight: 0.4,
				DirectoryWeight: 0.2,
				EnableColors:    true,
				ShellIntegration: "auto",
			}
			if err := saveConfig(configPath, cfg); err != nil {
				return cfg, nil // Return default config anyway
			}
			return cfg, nil
		}
		return nil, err
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func saveConfig(path string, cfg *Config) error {
	viper.Set("database_path", cfg.DatabasePath)
	viper.Set("max_suggestions", cfg.MaxSuggestions)
	viper.Set("recency_weight", cfg.RecencyWeight)
	viper.Set("frequency_weight", cfg.FrequencyWeight)
	viper.Set("directory_weight", cfg.DirectoryWeight)
	viper.Set("enable_colors", cfg.EnableColors)
	viper.Set("shell_integration", cfg.ShellIntegration)

	return viper.WriteConfigAs(path + ".yaml")
}

func Get() *Config {
	return cfg
}
