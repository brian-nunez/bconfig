package bconfig

import (
	"context"
)

// Config holds the merged configuration data.
type Config struct {
	data map[string]any
}

// New loads configuration from the given sources using context.Background.
func New(sources ...Source) (*Config, error) {
	return Load(context.Background(), sources...)
}

// Load loads configuration from the given sources in order.
// Later sources override values from earlier sources.
func Load(ctx context.Context, sources ...Source) (*Config, error) {
	return LoadWithOptions(ctx, sources)
}

// LoadWithOptions loads configuration from sources and applies options (e.g. validators).
func LoadWithOptions(ctx context.Context, sources []Source, opts ...Option) (*Config, error) {
	l := &loader{}
	for _, opt := range opts {
		opt(l)
	}

	data := make(map[string]any)
	for _, src := range sources {
		if src == nil {
			continue
		}
		srcData, err := src.Load(ctx)
		if err != nil {
			return nil, WrapSourceError(src.Name(), err)
		}
		data = deepMerge(data, srcData)
	}

	cfg := &Config{data: data}

	if l.validator != nil {
		if err := l.validator(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
