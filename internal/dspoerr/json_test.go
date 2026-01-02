package dspoerr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestMarshalJSON_Basic(t *testing.T) {
	err := New(CodeTimeout, "operation timed out")

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal: %v", marshalErr)
	}

	var result map[string]any
	if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal: %v", unmarshalErr)
	}

	if result["code"] != "TIMEOUT" {
		t.Errorf("expected code 'TIMEOUT', got %v", result["code"])
	}
	if result["message"] != "operation timed out" {
		t.Errorf("expected message 'operation timed out', got %v", result["message"])
	}

	// Hint and Details should not be present
	if _, hasHint := result["hint"]; hasHint {
		t.Error("expected hint to be omitted when empty")
	}
	if _, hasDetails := result["details"]; hasDetails {
		t.Error("expected details to be omitted when nil")
	}
}

func TestMarshalJSON_WithHint(t *testing.T) {
	err := New(CodeMissingPrereq, "prerequisite missing")
	err = WithHint(err, "Install the required package")

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal: %v", marshalErr)
	}

	var result map[string]any
	if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal: %v", unmarshalErr)
	}

	if result["code"] != "MISSING_PREREQ" {
		t.Errorf("expected code 'MISSING_PREREQ', got %v", result["code"])
	}
	if result["message"] != "prerequisite missing" {
		t.Errorf("expected message 'prerequisite missing', got %v", result["message"])
	}
	if result["hint"] != "Install the required package" {
		t.Errorf("expected hint 'Install the required package', got %v", result["hint"])
	}
}

func TestMarshalJSON_WithDetails(t *testing.T) {
	err := New(CodeUnsupportedPlatform, "platform not supported")
	err = WithDetails(err, map[string]any{
		"detected": "fedora-39",
		"count":    5,
		"enabled":  false,
	})

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal: %v", marshalErr)
	}

	var result map[string]any
	if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal: %v", unmarshalErr)
	}

	if result["code"] != "UNSUPPORTED_PLATFORM" {
		t.Errorf("expected code 'UNSUPPORTED_PLATFORM', got %v", result["code"])
	}

	details, ok := result["details"].(map[string]any)
	if !ok {
		t.Fatalf("expected details to be a map, got %T", result["details"])
	}
	if details["detected"] != "fedora-39" {
		t.Errorf("expected details.detected 'fedora-39', got %v", details["detected"])
	}
	if details["count"].(float64) != 5 {
		t.Errorf("expected details.count 5, got %v", details["count"])
	}
	if details["enabled"] != false {
		t.Errorf("expected details.enabled false, got %v", details["enabled"])
	}
}

func TestMarshalJSON_Complete(t *testing.T) {
	err := New(CodeStateCorrupt, "state file corrupted")
	err = WithHint(err, "Restore from backup or run dspo reset")
	err = WithDetails(err, map[string]any{
		"file": testStatePath,
		"size": 1024,
	})

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal: %v", marshalErr)
	}

	var result map[string]any
	if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal: %v", unmarshalErr)
	}

	if result["code"] != "STATE_CORRUPT" {
		t.Errorf("expected code 'STATE_CORRUPT', got %v", result["code"])
	}
	if result["message"] != "state file corrupted" {
		t.Errorf("expected message 'state file corrupted', got %v", result["message"])
	}
	if result["hint"] != "Restore from backup or run dspo reset" {
		t.Errorf("expected hint 'Restore from backup or run dspo reset', got %v", result["hint"])
	}

	details, ok := result["details"].(map[string]any)
	if !ok {
		t.Fatalf("expected details to be a map")
	}
	if details["file"] != testStatePath {
		t.Errorf("expected details.file '%s', got %v", testStatePath, details["file"])
	}
}

func TestMarshalJSON_CauseOmitted(t *testing.T) {
	cause := fmt.Errorf("internal system error with sensitive data")
	err := Wrap(CodeStateWriteFailed, "failed to write state", cause)

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal: %v", marshalErr)
	}

	var result map[string]any
	if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal: %v", unmarshalErr)
	}

	// CRITICAL: Cause should NOT be marshaled
	if _, hasCause := result["cause"]; hasCause {
		t.Error("Cause should NOT be marshaled to JSON (internal-only per CCC)")
	}

	if result["code"] != "STATE_WRITE_FAILED" {
		t.Errorf("expected code 'STATE_WRITE_FAILED', got %v", result["code"])
	}
	if result["message"] != "failed to write state" {
		t.Errorf("expected message 'failed to write state', got %v", result["message"])
	}

	// Verify Cause is still accessible via Unwrap
	if err.Unwrap() != cause {
		t.Error("expected Unwrap to still return the cause")
	}
}

func TestMarshalJSON_ExpectedEnvelope(t *testing.T) {
	// Test exact format from CCC §1 (lines 108-118)
	err := New(CodeUnsupportedPlatform, "Fedora 39 is not supported (supported: current and N-1).")
	err = WithHint(err, "Upgrade to Fedora current or run `dspo doctor --json` to see supported targets.")
	err = WithDetails(err, map[string]any{"detected": "fedora-39"})

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal: %v", marshalErr)
	}

	var actual map[string]any
	if unmarshalErr := json.Unmarshal(data, &actual); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal: %v", unmarshalErr)
	}

	expected := map[string]any{
		"code":    "UNSUPPORTED_PLATFORM",
		"message": "Fedora 39 is not supported (supported: current and N-1).",
		"hint":    "Upgrade to Fedora current or run `dspo doctor --json` to see supported targets.",
		"details": map[string]any{"detected": "fedora-39"},
	}

	// Compare code
	if actual["code"] != expected["code"] {
		t.Errorf("code mismatch: expected %v, got %v", expected["code"], actual["code"])
	}

	// Compare message
	if actual["message"] != expected["message"] {
		t.Errorf("message mismatch: expected %v, got %v", expected["message"], actual["message"])
	}

	// Compare hint
	if actual["hint"] != expected["hint"] {
		t.Errorf("hint mismatch: expected %v, got %v", expected["hint"], actual["hint"])
	}

	// Compare details
	actualDetails, ok := actual["details"].(map[string]any)
	if !ok {
		t.Fatalf("expected details to be a map, got %T", actual["details"])
	}
	expectedDetails := expected["details"].(map[string]any)
	if !reflect.DeepEqual(actualDetails, expectedDetails) {
		t.Errorf("details mismatch: expected %v, got %v", expectedDetails, actualDetails)
	}

	// Verify no extra fields
	if len(actual) != 4 {
		t.Errorf("expected exactly 4 fields in JSON, got %d: %v", len(actual), actual)
	}
}
