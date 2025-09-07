package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultFormat  string                 `yaml:"default_format"`
	Quality        uint16                 `yaml:"quality"`
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
	// Batch processing options
	BatchProcessing BatchConfig `yaml:"batch_processing"`
}

// BatchConfig contains configuration for batch processing features
type BatchConfig struct {
	RecursiveSearch   bool   `yaml:"recursive_search"`   // Search subdirectories recursively
	MaxDepth          int    `yaml:"max_depth"`          // Maximum directory depth to search (0 = unlimited)
	PreserveStructure bool   `yaml:"preserve_structure"` // Preserve directory structure in output
	OutputDir         string `yaml:"output_dir"`         // Custom output directory for batch processing
	GroupByFolder     bool   `yaml:"group_by_folder"`    // Group results by source folder
	SkipEmptyDirs     bool   `yaml:"skip_empty_dirs"`    // Skip directories with no images
	FollowSymlinks    bool   `yaml:"follow_symlinks"`    // Follow symbolic links
}

// DefaultConfig returns the default configuration for gopix.
// The returned configuration is a reasonable set of defaults, but can be overridden
// by the user through the command line flags or a configuration file.
//
// The default configuration is as follows:
//
// - Default format: png
// - Quality: 80
// - Number of workers: the number of CPUs available
// - Maximum dimension: 0 (no limit)
// - Log level: info
// - Supported extentions: png, jpg, jpeg, webp
// - Auto backup: true
// - Resume enabled: true
// - Keep original: false
// - Dry run: false
// - Verbose logging: false
//
// The output settings are as follows:
//
// - For PNG: use best speed compression
// - For JPG: use quality 80
// - For JPEG: use quality 80
// - For WebP: use quality 80 and lossless compression
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
		BatchProcessing: BatchConfig{
			RecursiveSearch:   true,
			MaxDepth:          0, // 0 = unlimited depth
			PreserveStructure: true,
			OutputDir:         "",
			GroupByFolder:     false,
			SkipEmptyDirs:     true,
			FollowSymlinks:    false,
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

	//load existing config - use ReadFile for better performance
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

// Save writes the current configuration to a YAML file in the user's config directory.
// It marshals the Config struct to YAML format and saves it as "config.yaml".
// If the marshaling or file writing fails, it returns an error detailing the failure.

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
