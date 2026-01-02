package dspoerr

import (
	"errors"
	"fmt"
	"testing"
)

const testStatePath = "/path/to/state.json"

func TestNew(t *testing.T) {
	err := New(CodeStepFailed, "test error")

	if err.Code != CodeStepFailed {
		t.Errorf("expected code %s, got %s", CodeStepFailed, err.Code)
	}
	if err.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", err.Message)
	}
	if err.Cause != nil {
		t.Error("expected Cause to be nil")
	}
	if err.Hint != "" {
		t.Errorf("expected Hint to be empty, got '%s'", err.Hint)
	}
	if err.Details != nil {
		t.Errorf("expected Details to be nil, got %v", err.Details)
	}

	// Test Error() format
	expected := "STEP_FAILED: test error"
	if err.Error() != expected {
		t.Errorf("expected Error() '%s', got '%s'", expected, err.Error())
	}
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("underlying error")
	err := Wrap(CodeStepFailed, "step failed", cause)

	if err.Code != CodeStepFailed {
		t.Errorf("expected code %s, got %s", CodeStepFailed, err.Code)
	}
	if err.Message != "step failed" {
		t.Errorf("expected message 'step failed', got '%s'", err.Message)
	}
	if err.Cause != cause {
		t.Error("expected cause to be preserved")
	}
	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to work via Unwrap")
	}

	// Test Error() format
	expected := "STEP_FAILED: step failed"
	if err.Error() != expected {
		t.Errorf("expected Error() '%s', got '%s'", expected, err.Error())
	}
}

func TestWithHint(t *testing.T) {
	err := New(CodeMissingPrereq, "missing prerequisite")
	result := WithHint(err, "Install the missing package")

	if result != err {
		t.Error("expected WithHint to return the same error for chaining")
	}
	if err.Hint != "Install the missing package" {
		t.Errorf("expected hint 'Install the missing package', got '%s'", err.Hint)
	}
}

func TestWithDetails(t *testing.T) {
	err := New(CodeUnsupportedPlatform, "unsupported platform")
	details := map[string]any{
		"detected": "fedora-39",
		"count":    42,
		"enabled":  true,
	}
	result := WithDetails(err, details)

	if result != err {
		t.Error("expected WithDetails to return the same error for chaining")
	}
	if err.Details == nil {
		t.Fatal("expected Details to be set")
	}
	if err.Details["detected"] != "fedora-39" {
		t.Errorf("expected Details[detected] 'fedora-39', got %v", err.Details["detected"])
	}
	if err.Details["count"] != 42 {
		t.Errorf("expected Details[count] 42, got %v", err.Details["count"])
	}
	if err.Details["enabled"] != true {
		t.Errorf("expected Details[enabled] true, got %v", err.Details["enabled"])
	}
}

func TestEnricherChaining(t *testing.T) {
	err := New(CodeStateCorrupt, "state corrupted")
	details := map[string]any{"file": testStatePath}

	result := WithDetails(WithHint(err, "Run dspo doctor"), details)

	if result != err {
		t.Error("expected chaining to return the same error")
	}
	if err.Hint != "Run dspo doctor" {
		t.Errorf("expected hint 'Run dspo doctor', got '%s'", err.Hint)
	}
	if err.Details == nil || err.Details["file"] != testStatePath {
		t.Errorf("expected Details[file] '%s', got %v", testStatePath, err.Details)
	}

	// Test reverse order
	err2 := New(CodeLockHeld, "lock held")
	result2 := WithHint(WithDetails(err2, details), "Release the lock")

	if result2 != err2 {
		t.Error("expected reverse chaining to return the same error")
	}
	if err2.Hint != "Release the lock" {
		t.Errorf("expected hint 'Release the lock', got '%s'", err2.Hint)
	}
	if err2.Details == nil || err2.Details["file"] != testStatePath {
		t.Errorf("expected Details[file] '%s', got %v", testStatePath, err2.Details)
	}
}

func TestAs(t *testing.T) {
	// Test with dspoerr.Error
	dspoErr := New(CodeTimeout, "timeout occurred")
	result, ok := As(dspoErr)
	if !ok {
		t.Error("expected As to return true for dspoerr.Error")
	}
	if result != dspoErr {
		t.Error("expected As to return the same error")
	}

	// Test with wrapped dspoerr.Error
	wrappedErr := fmt.Errorf("wrapper: %w", dspoErr)
	result2, ok2 := As(wrappedErr)
	if !ok2 {
		t.Error("expected As to return true for wrapped dspoerr.Error")
	}
	if result2 != dspoErr {
		t.Error("expected As to extract the dspoerr.Error from chain")
	}

	// Test with non-dspoerr error
	stdErr := errors.New("standard error")
	result3, ok3 := As(stdErr)
	if ok3 {
		t.Error("expected As to return false for non-dspoerr error")
	}
	if result3 != nil {
		t.Error("expected As to return nil for non-dspoerr error")
	}

	// Test with nil error
	result4, ok4 := As(nil)
	if ok4 {
		t.Error("expected As to return false for nil error")
	}
	if result4 != nil {
		t.Error("expected As to return nil for nil error")
	}
}

func TestErrorChainPreservation(t *testing.T) {
	// Create error chain: standard error -> wrapped -> wrapped again
	baseErr := errors.New("base error")
	err1 := Wrap(CodeStepFailed, "step 1 failed", baseErr)
	err2 := Wrap(CodeStateWriteFailed, "state write failed", err1)

	// Test errors.Is works through chain
	if !errors.Is(err2, baseErr) {
		t.Error("expected errors.Is to work through error chain")
	}
	if !errors.Is(err2, err1) {
		t.Error("expected errors.Is to find intermediate error")
	}

	// Test errors.As works through chain
	var dspoErr *Error
	if !errors.As(err2, &dspoErr) {
		t.Fatal("expected errors.As to extract dspoerr.Error from chain")
	}
	if dspoErr.Code != CodeStateWriteFailed {
		t.Errorf("expected errors.As to extract outermost error with code %s, got %s",
			CodeStateWriteFailed, dspoErr.Code)
	}

	// Test As helper works
	extracted, ok := As(err2)
	if !ok {
		t.Fatal("expected As helper to work on error chain")
	}
	if extracted.Code != CodeStateWriteFailed {
		t.Errorf("expected As to extract code %s, got %s",
			CodeStateWriteFailed, extracted.Code)
	}
}

func TestUnwrap(t *testing.T) {
	// Test error with cause
	cause := errors.New("underlying cause")
	err := Wrap(CodeNetworkRequired, "network required", cause)

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Error("expected Unwrap to return the cause")
	}

	// Test error without cause
	err2 := New(CodeSudoRequired, "sudo required")
	unwrapped2 := err2.Unwrap()
	if unwrapped2 != nil {
		t.Error("expected Unwrap to return nil for error without cause")
	}
}
