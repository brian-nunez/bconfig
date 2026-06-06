package env

const DriverName = "env"

type Config struct {
	// Prefixes specifies the prefixes to load from environment variables.
	Prefixes []string
}

// DriverName returns the name of this driver.
func (Config) DriverName() string {
	return DriverName
}
