package file

const DriverName = "file"

type Config struct {
	// Path is the path to the configuration file.
	Path string
}

// DriverName returns the name of this driver.
func (Config) DriverName() string {
	return DriverName
}
