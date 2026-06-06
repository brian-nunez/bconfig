package bconfig

type loader struct {
	validator  func(*Config) error
	strictMode bool
}

// Option is a function that configures the loader.
type Option func(*loader)

// WithValidator sets a validation function to be run after configuration is loaded.
func WithValidator(fn func(*Config) error) Option {
	return func(l *loader) {
		l.validator = fn
	}
}

// WithStrictMode enables or disables strict mode.
func WithStrictMode(enabled bool) Option {
	return func(l *loader) {
		l.strictMode = enabled
	}
}
