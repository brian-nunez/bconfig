package bconfig

import (
	"fmt"
	"strings"
)

// Get retrieves the value at the given dot-path.
// Returns nil if the path is not found.
func (c *Config) Get(path string) any {
	if path == "" {
		return nil
	}
	parts := strings.Split(path, ".")
	var current any = c.data
	for _, part := range parts {
		m, ok := toMap(current)
		if !ok {
			return nil
		}
		val, exists := m[part]
		if !exists {
			return nil
		}
		current = val
	}
	return current
}

// String retrieves the value at the given dot-path as a string.
// Returns a empty string if the value is missing.
func (c *Config) String(path string) string {
	val := c.Get(path)
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}

// Bool retrieves the value at the given dot-path as a bool.
// Returns false if the value is missing or not a bool.
func (c *Config) Bool(path string) bool {
	val := c.Get(path)
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// Int retrieves the value at the given dot-path as an int.
// Returns 0 if the value is missing or not a number.
func (c *Config) Int(path string) int {
	val := c.Get(path)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int64:
		return int(v)
	case int32:
		return int(v)
	}
	return 0
}

// Float retrieves the value at the given dot-path as a float64.
// Returns 0.0 if the value is missing or not a number.
func (c *Config) Float(path string) float64 {
	val := c.Get(path)
	if val == nil {
		return 0.0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	}
	return 0.0
}

// Data returns a deep copy of the configuration data map.
func (c *Config) Data() map[string]any {
	return deepMerge(nil, c.data)
}
