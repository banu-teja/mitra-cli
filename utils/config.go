package utils

import (
	"embed"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type ModelConfig struct {
	Name    string `mapstructure:"NAME"`
	BaseURL string `mapstructure:"BASE_URL"`
	APIKey  string `mapstructure:"API_KEY"`
}

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	DBSource    string `mapstructure:"DB_SOURCE"`

	// Map of model configurations
	Models map[string]ModelConfig `mapstructure:"AI_MODELS"`
}

//go:embed app.env
var embeddedConfig embed.FS

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// Important:  We *don't* need AddConfigPath when using embedded files.
	// viper.AddConfigPath(path)  // Remove this line

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the embedded config file
	content, err := embeddedConfig.ReadFile("app.env")
	if err != nil {
		return config, fmt.Errorf("failed to read embedded config: %w", err)
	}

	// Use ReadConfig with a strings.Reader
	reader := strings.NewReader(string(content))
	err = viper.ReadConfig(reader)
	if err != nil {
		return config, fmt.Errorf("failed to parse embedded config: %w", err)
	}

	// Unmarshal the configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
