package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	DefaultFormat  string                 `yaml:"default_format"`
	Quality        uint16                  `yaml:"quality"`
	Workers        uint8                  `yaml:"workers"`
	MaxDimension   uint16                 `yaml:"max_dimension"`
	LogLevel       string                 `yaml:"log_level"`
	Extentions     []string               `yaml:"extentions"`
	OutputSettings map[string]interface{} `yaml:"output_settings"`
	AutoBackup     bool                   `yaml:"auto_backup"`
	ResumeEnabled  bool                   `yaml:"resume_enabled"`
	KeepOriginal   bool                   `yaml:"keep_original"`
	DryRun         bool                   `yaml:"dry_run"`
	Verbose        bool                   `yaml:"verbose"`
}

func DefaultConfig() *Config {
	return &Config{
		DefaultFormat: "png",
		Quality:       80,
		Workers:       uint8(runtime.NumCPU()),
		MaxDimension:  0,
		LogLevel:      "info",
		Extentions:    []string{"png", "jpg", "jpeg", "webp"},
		AutoBackup:    true,
		ResumeEnabled: true,
		KeepOriginal:  false,
		DryRun:        false,
	    // Verbose:       false,
		OutputSettings: map[string]interface{}{
			"png": map[string]interface{}{
				"compression": "best_speed",
			},
			"jpg": map[string]interface{}{
				"quality": 80,
			},
			"jpeg": map[string]interface{}{
				"quality": 80,
			},
			"webp": map[string]interface{}{
				"quality":  80,
				"lossless": false,
			},
		},
	}
}

func getConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".gopix/config")
}

func LoadConfig() (*Config, error) {
	configDirectory := getConfigDirectory()
	configPath := filepath.Join(configDirectory, "config.yaml")

	//create config directory if it doesn't exist
	if err := os.MkdirAll(configDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	// if config file doesn't exist, create it
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := DefaultConfig()
		if err := defaultConfig.Save(); err != nil {
			return nil, fmt.Errorf("failed to save default config: %v", err)
		}
		return defaultConfig, nil
	}

	//load existing config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}
	return &conf, nil
}

func (c *Config) Save() error {
	configDirctory := getConfigDirectory()
	configPath := filepath.Join(configDirctory, "config.yaml")
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	return nil
}
