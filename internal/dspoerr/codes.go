package dspoerr

import "fmt"

// ExitCode represents process exit codes
type ExitCode int

const (
	// ExitSuccess indicates successful execution
	ExitSuccess ExitCode = 0
	// ExitGenericError indicates a general error
	ExitGenericError ExitCode = 1
	// ExitValidationError indicates input validation failure
	ExitValidationError ExitCode = 2
	// ExitPlatformError indicates platform incompatibility
	ExitPlatformError ExitCode = 3
)

// Error represents a dspo error with an exit code
type Error struct {
	Code    ExitCode
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// New creates a new Error
func New(code ExitCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with context and exit code
func Wrap(code ExitCode, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
