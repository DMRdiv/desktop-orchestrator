package dspoerr

import (
	"encoding/json"
	"errors"
	"fmt"
)

//
// Error type
//

// Error represents a structured, user-facing error produced by dspo.
//
// It implements the standard `error` interface.
// The Cause field is internal-only and must never be printed in
// non-debug user output.
type Error struct {
	Code    Code           `json:"code"`
	Message string         `json:"message"`
	Hint    string         `json:"hint,omitempty"`
	Details map[string]any `json:"details,omitempty"`
	Cause   error          `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// MarshalJSON implements custom JSON marshaling
func (e *Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

//
// Constructors
//

// New creates a new dspo error with the given code and message.
func New(code Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// Wrap creates a new dspo error that wraps an underlying cause.
func Wrap(code Code, msg string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}

//
// Optional enrichers
//

// WithHint attaches a user-facing remediation hint.
func WithHint(e *Error, hint string) *Error {
	e.Hint = hint
	return e
}

// WithDetails attaches structured diagnostic details.
// Intended for debug output only.
func WithDetails(e *Error, details map[string]any) *Error {
	e.Details = details
	return e
}

//
// Type assertion helper
//

// As attempts to extract a *dspoerr.Error from an error chain.
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}
