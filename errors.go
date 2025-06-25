package gogacon

import "fmt"

// ConfigError represents configuration-related errors with context
type ConfigError struct {
	Operation string
	Path      string
	Err       error
}

// Error - formats the error message with contextual information.
func (e ConfigError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("config error: %s: %s", e.Operation, e.Err)
	}
	return fmt.Sprintf("config error: %s %q: %v", e.Operation, e.Path, e.Err)
}

func (e ConfigError) Unwrap() error { return e.Err }

// NewError creates new ConfigError with contex
//
// args:
//
//	operation: running operation when error occured
//	path: path to configuration file (optional)
//	err: initial error
func NewError(operation string, path string, err error) ConfigError {
	return ConfigError{
		Operation: operation,
		Path:      path,
		Err:       err,
	}
}
