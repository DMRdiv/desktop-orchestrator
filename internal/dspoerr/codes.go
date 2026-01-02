package dspoerr

//
// Error codes
//

// Code represents a SCREAMING_SNAKE_CASE error code identifier
type Code string

// Canonical error codes
const (
	CodeUnsupportedPlatform   Code = "UNSUPPORTED_PLATFORM"
	CodeMissingPrereq         Code = "MISSING_PREREQ"
	CodeInvalidProfileSchema  Code = "INVALID_PROFILE_SCHEMA"
	CodeMergeConflict         Code = "MERGE_CONFLICT"
	CodeLockHeld              Code = "LOCK_HELD"
	CodeSudoRequired          Code = "SUDO_REQUIRED"
	CodeStepFailed            Code = "STEP_FAILED"
	CodeStateCorrupt          Code = "STATE_CORRUPT"
	CodeStateWriteFailed      Code = "STATE_WRITE_FAILED"
	CodeAssetChecksumMissing  Code = "ASSET_CHECKSUM_MISSING"
	CodeAssetChecksumMismatch Code = "ASSET_CHECKSUM_MISMATCH"
	CodeNetworkRequired       Code = "NETWORK_REQUIRED"
	CodeTimeout               Code = "TIMEOUT"
)
