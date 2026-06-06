package bconfig

import (
	"fmt"
)

// WrapSourceError wraps a source error with the expected format:
// "bconfig: failed to load source <source_name>: <original_error>"
func WrapSourceError(sourceName string, err error) error {
	return fmt.Errorf("bconfig: failed to load source %s: %w", sourceName, err)
}

// WrapDecodeError wraps a decode error with the expected format:
// "bconfig: decode failed: <original_error>"
func WrapDecodeError(err error) error {
	return fmt.Errorf("bconfig: decode failed: %w", err)
}
