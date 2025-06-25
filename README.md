# GoGACon - Go Configuration Manager Library

[![Go Reference](https://pkg.go.dev/badge/github.com/Galdoba/gogacon.svg)](https://pkg.go.dev/github.com/Galdoba/gogacon)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GoGACon is a lightweight, zero-dependency Go library designed to simplify configuration management in Go applications. It provides a straightforward way to handle configuration loading, default value initialization, and error handling with rich context.

## Features

- 🛠 Automatic Config Initialization: Creates default config files on first run
- 📁 Standard Path Handling: Uses XDG-compliant paths (~/.config/<appname>/)
- 🧩 Format Agnostic: Works with any format via the Serializer interface
- 🚨 Rich Error Context: Detailed error messages with operation and path information
- 🛡 Validation: Ensures required parameters are provided
- 📂 Directory Management: Automatically creates necessary directories

## Installation

go get github.com/Galdoba/gogacon

## Quick Start

### 1. Implement Serializer Interface

    package main

    import "gogacon"

    type AppConfig struct {
        Port     int    `json:"port"`
        LogLevel string `json:"log_level"`
    }

    func (c *AppConfig) Marshal() ([]byte, error) {
        // Implement your serialization logic (JSON, YAML, etc.)
    }

    func (c *AppConfig) Unmarshal(data []byte) error {
        // Implement your deserialization logic
    }

### 2. Initialize and Use Config Manager

    package main

    import (
    "fmt"
    "gogacon"
    )

    func main() {
        defaults := gogacon.Defaults{
            AppName: "myapp",
            DefaultConfigValues: &AppConfig{
                Port:     8080,
                LogLevel: "info",
            },
        }

        manager, err := gogacon.NewConfigManager(defaults)
        if err != nil {
            panic(fmt.Errorf("config initialization failed: %w", err))
        }

        var config AppConfig
        if err := manager.LoadConfig("", &config); err != nil {
            panic(fmt.Errorf("config loading failed: %w", err))
        }

        fmt.Printf("Server running on port %d\n", config.Port)
    }

## Core Concepts

### ConfigManager
The main component that handles configuration loading and initialization.

Methods:
- NewConfigManager(defaults Defaults) (*ConfigManager, error)  
  Creates a new configuration manager with validation
- LoadConfig(filePath string, target Serializer) error  
  Loads configuration from specified path or default location

### Serializer Interface
go
type Serializer interface {
    Marshal() ([]byte, error)
    Unmarshal(data []byte) error
}
Implement this interface to support any configuration format (JSON, YAML, TOML, etc.)

### Error Handling
Errors are returned as ConfigError type with rich context:
go
type ConfigError struct {
    Operation string // Operation being performed
    Path      string // Configuration file path
    Err       error  // Underlying error
}

Example error message:  
config error: unmarshal config "/home/user/.config/myapp/default.conf": invalid character '}'

## Advanced Usage

### Custom Configuration Path
go
// Load from specific location
err := manager.LoadConfig("/etc/myapp/special.conf", &config)

### Handling Errors
go
if err := manager.LoadConfig("", &config); err != nil {
    var cfgErr gogacon.ConfigError
    if errors.As(err, &cfgErr) {
        fmt.Printf("Operation: %s\nPath: %s\nError: %v\n", 
            cfgErr.Operation, cfgErr.Path, cfgErr.Err)
    } else {
        fmt.Printf("Unexpected error: %v\n", err)
    }
    os.Exit(1)
}

### Implementing Custom Serializers
go
// JSON Serializer Example
type JSONConfig struct {
    Data map[string]interface{}
}

func (j *JSONConfig) Marshal() ([]byte, error) {
    return json.MarshalIndent(j.Data, "", "  ")
}

func (j *JSONConfig) Unmarshal(data []byte) error {
    return json.Unmarshal(data, &j.Data)
}


## Error Reference

| Error Condition | Example Error Message |
|----------------|------------------------|
| Missing AppName | AppName must be specified |
| Missing DefaultConfigValues | DefaultConfigValues must be specified |
| Permission Denied | config error: create config directory "/path": permission denied |
| Invalid Config | config error: unmarshal config "/path": invalid syntax |
| File Not Found | config error: read config "/path": file not found |
| Home Dir Unavailable | config error: resolve config path: failed to get home dir |

## Contributing

Contributions are welcome! Please follow these steps:
1. Open an issue describing the problem or feature
2. Fork the repository and create your feature branch
3. Add tests for your changes
4. Submit a pull request with detailed description

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.