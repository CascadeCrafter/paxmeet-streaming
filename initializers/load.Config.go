package initializers

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the top-level structure of the YAML configuration
type Config struct {
	Auth    AuthConfig    `yaml:"auth"`
	LiveKit LiveKitConfig `yaml:"livekit"`
	Redis   RedisConfig   `yaml:"redis"`
	Backend BackendConfig `yaml:"backend"`
}

// LiveKitConfig represents the nested "livekit" structure in the YAML
type LiveKitConfig struct {
	Uri       string `yaml:"uri"`
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

// AuthConfig represents the nested "auth" structure in the YAML
type AuthConfig struct {
	Uri string `yaml:"uri"`
}

// BackendConfig represents the nested "auth" structure in the YAML
type BackendConfig struct {
	Uri string `yaml:"uri"`
}

// RedisConfig represents the nested "redis" structure in the YAML
type RedisConfig struct {
	Url string `yaml:"url"`
}

// ReadConfig reads and unmarshals YAML configuration from a file
func LoadConfig() (*Config, error) {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
