package bconfig

import "context"

// Source represents a configuration source that can load data.
type Source interface {
	Name() string
	Load(ctx context.Context) (map[string]any, error)
}
