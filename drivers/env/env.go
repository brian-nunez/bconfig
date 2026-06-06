package env

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/brian-nunez/bconfig"
)

type envSource struct {
	prefixes []string
}

// Source returns a bconfig.Source that loads configuration from environment variables.
// If prefixes are provided, only variables starting with one of the prefixes are loaded,
// and the prefix is stripped before mapping.
func Source(prefix ...string) bconfig.Source {
	return &envSource{prefixes: prefix}
}

// Name returns the name of this source.
func (e *envSource) Name() string {
	if len(e.prefixes) == 0 {
		return "env"
	}
	return "env:" + strings.Join(e.prefixes, ",")
}

// Load loads and maps environment variables.
func (e *envSource) Load(ctx context.Context) (map[string]any, error) {
	result := make(map[string]any)
	environ := os.Environ()

	for _, envVar := range environ {
		pair := strings.SplitN(envVar, "=", 2)
		key := pair[0]
		val := pair[1]

		if len(e.prefixes) > 0 {
			matched := false
			for _, p := range e.prefixes {
				if strings.HasPrefix(key, p) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// Lowercase the key and map it directly
		keyLower := strings.ToLower(key)
		result[keyLower] = parseValue(val)
	}

	return result, nil
}

// parseValue converts environment value strings into basic types.
// `"true"` -> bool true
// `"false"` -> bool false
// integer strings -> int
// float strings -> float64
// all others -> string
func parseValue(val string) any {
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	return val
}
