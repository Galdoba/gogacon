package gogacon_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Galdoba/gogacon"
)

// MockSerializer реализует интерфейс Serializer для тестирования
type MockSerializer struct {
	MarshalData    []byte
	MarshalErr     error
	UnmarshalErr   error
	UnmarshalCalls int
}

func (m *MockSerializer) Marshal() ([]byte, error) {
	return m.MarshalData, m.MarshalErr
}

func (m *MockSerializer) Unmarshal(data []byte) error {
	m.UnmarshalCalls++
	return m.UnmarshalErr
}

func TestConfigManager_FirstRun(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}
	// Создаем временную директорию
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Мок сериализатора с дефолтными значениями
	mockSerializer := &MockSerializer{
		MarshalData: []byte("config data"),
	}

	manager, err := gogacon.NewConfigManager(gogacon.Defaults{
		AppName:             "testapp",
		DefaultConfigValues: mockSerializer,
	})
	if err != nil {
		t.Fatalf("NewConfigManager failed: %v", err)
	}

	// Загружаем конфиг (должен создать файл)
	target := &MockSerializer{}
	if err := manager.LoadConfig("", target); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Проверяем что файл создан
	configPath := filepath.Join(tempDir, ".config", "testapp", "default.conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file not created")
	}

	// Проверяем что дефолтные значения были записаны
	content, _ := os.ReadFile(configPath)
	if string(content) != "config data" {
		t.Errorf("Invalid config content: %s", content)
	}
}

func TestLoadConfig_ExistingConfig(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Создаем предварительно заполненный конфиг
	configPath := filepath.Join(tempDir, ".config", "testapp", "default.conf")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	os.WriteFile(configPath, []byte("existing config"), 0644)

	manager, _ := gogacon.NewConfigManager(gogacon.Defaults{
		AppName:             "testapp",
		DefaultConfigValues: &MockSerializer{},
	})

	target := &MockSerializer{}
	if err := manager.LoadConfig("", target); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Проверяем что конфиг был загружен
	if target.UnmarshalCalls != 1 {
		t.Error("Unmarshal not called")
	}
}

func TestLoadConfig_InvalidConfig(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Создаем поврежденный конфиг
	configPath := filepath.Join(tempDir, ".config", "testapp", "default.conf")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	os.WriteFile(configPath, []byte("invalid data"), 0644)

	manager, _ := gogacon.NewConfigManager(gogacon.Defaults{
		AppName: "testapp",
		DefaultConfigValues: &MockSerializer{
			MarshalData: []byte("default data"),
		},
	})

	target := &MockSerializer{
		UnmarshalErr: errors.New("parse error"),
	}

	err := manager.LoadConfig("", target)
	if err == nil {
		t.Fatal("Expected error for invalid config")
	}

	// Проверяем тип ошибки
	if _, ok := err.(gogacon.ConfigError); !ok {
		t.Errorf("Expected ConfigError, got %T", err)
	}

	// Проверяем содержание ошибки
	expectedMsg := `config error: unmarshal config "` + configPath + `": parse error`
	if err.Error() != expectedMsg {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestLoadConfig_PermissionDenied(t *testing.T) {
	// Пропускаем тест на Windows из-за разных моделей прав доступа
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Создаем директорию без прав на запись
	configDir := filepath.Join(tempDir, ".config", "testapp")
	os.MkdirAll(configDir, 0555) // Только чтение

	manager, _ := gogacon.NewConfigManager(gogacon.Defaults{
		AppName: "testapp",
		DefaultConfigValues: &MockSerializer{
			MarshalData: []byte("test data"),
		},
	})

	err := manager.LoadConfig("", &MockSerializer{})
	if err == nil {
		t.Fatal("Expected permission error")
	}

	// Проверяем содержание ошибки
	expectedOp := "create config directory"
	if err.(gogacon.ConfigError).Operation != expectedOp {
		t.Errorf("Expected operation %q, got %q", expectedOp, err.(gogacon.ConfigError).Operation)
	}
}

func TestLoadConfig_SpecificPath(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}
	tempDir := t.TempDir()
	customPath := filepath.Join(tempDir, "custom.conf")
	os.WriteFile(customPath, []byte("custom config"), 0644)

	manager, _ := gogacon.NewConfigManager(gogacon.Defaults{
		AppName:             "testapp",
		DefaultConfigValues: &MockSerializer{},
	})

	target := &MockSerializer{}
	if err := manager.LoadConfig(customPath, target); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Проверяем что загрузка произошла из правильного файла
	if target.UnmarshalCalls != 1 {
		t.Error("Config not loaded from custom path")
	}
}

func TestNewConfigManager_Validation(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}
	tests := []struct {
		name     string
		defaults gogacon.Defaults
		errMsg   string
	}{
		{
			name:     "Missing AppName",
			defaults: gogacon.Defaults{DefaultConfigValues: &MockSerializer{}},
			errMsg:   "AppName must be specified",
		},
		{
			name:     "Missing DefaultConfigValues",
			defaults: gogacon.Defaults{AppName: "testapp"},
			errMsg:   "DefaultConfigValues must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gogacon.NewConfigManager(tt.defaults)
			if err == nil || err.Error() != tt.errMsg {
				t.Errorf("Expected error %q, got %v", tt.errMsg, err)
			}
		})
	}
}
