package dspoerr

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ExitValidationError, "test error")
	if err.Code != ExitValidationError {
		t.Errorf("expected code %d, got %d", ExitValidationError, err.Code)
	}
	if err.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", err.Message)
	}
	if err.Error() != "test error" {
		t.Errorf("expected Error() 'test error', got '%s'", err.Error())
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrap(ExitPlatformError, "wrapped", cause)

	if err.Code != ExitPlatformError {
		t.Errorf("expected code %d, got %d", ExitPlatformError, err.Code)
	}
	if err.Cause != cause {
		t.Error("expected cause to be preserved")
	}
	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to work via Unwrap")
	}
}
