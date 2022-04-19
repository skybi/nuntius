package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"strings"
)

// Config represents the application configuration structure
type Config struct {
	Environment string `default:"prod"`

	APIAddress string `default:"http://localhost:8082" split_words:"true"`
	APIKey     string `split_words:"true"`

	FeedMETARs bool `envconfig:"feed_metars"`
}

// LoadFromEnv loads a new configuration structure using environment variables and an optional .env file
func LoadFromEnv() (*Config, error) {
	// Load a .env file if it exists
	_ = godotenv.Overload()

	// Load a new configuration structure using environment variables
	config := new(Config)
	if err := envconfig.Process("sbf", config); err != nil {
		return nil, err
	}
	return config, nil
}

// IsEnvProduction returns whether the application runs in production environment
func (config *Config) IsEnvProduction() bool {
	return strings.ToLower(config.Environment) != "dev"
}
