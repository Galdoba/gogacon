package gogacon

import (
	"fmt"
	"os"
	"path/filepath"
)

// Defaults contains default data for ConfigManager
type Defaults struct {
	AppName             string     //Application Name (required)
	DefaultConfigValues Serializer //Default configuration values (required)
}

// ConfigManager manages loading and savings of configuration
type ConfigManager struct {
	defaults Defaults
	filePath string
}

// NewConfigManager creates new ConfigManager instance
// Checks required fields of Defaults
func NewConfigManager(d Defaults) (*ConfigManager, error) {
	if d.AppName == "" {
		return nil, NewError("initialization", "", fmt.Errorf("AppName must be specified"))
	}
	if d.DefaultConfigValues == nil {
		return nil, NewError("initialization", "", fmt.Errorf("DefaultConfigValues must be specified"))
	}
	return &ConfigManager{defaults: d}, nil
}

// buildPath creates default filepath to configuration file
// Returns absolute path or error
func (cm *ConfigManager) buildPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", cm.defaults.AppName, "default.conf"), nil
}

func (cm *ConfigManager) ensureConfigFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return NewError("create config directory", dir, err)
	}

	data, err := cm.defaults.DefaultConfigValues.Marshal()
	if err != nil {
		return NewError("marshal default config", path, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return NewError("create default config", path, err)
	}
	return nil
}

// LoadConfig loads configuration from path provided or from default path
// Creates configuration file with default values at default filepath
func (cm *ConfigManager) LoadConfig(filePath string, target Serializer) error {
	if filePath == "" {
		return fmt.Errorf("no path provided")
	}
	if err := cm.ensureConfigFile(filePath); err != nil {
		return err
	}

	bt, err := os.ReadFile(filePath)
	if err != nil {
		return NewError("read config", filePath, err)
	}

	if err = target.Unmarshal(bt); err != nil {
		return NewError("unmarshal config", filePath, err)
	}
	cm.filePath = filePath
	return nil
}

// SaveConfig saves configuration to path it was loaded from
func (cm *ConfigManager) SaveConfig(config Serializer) error {
	bt, err := config.Marshal()
	if err != nil {
		return NewError("marshal config", "", err)
	}
	if err := os.WriteFile(cm.filePath, bt, 0644); err != nil {
		return NewError("save config", cm.filePath, err)
	}
	return nil
}
