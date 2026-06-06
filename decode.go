package bconfig

import (
	"encoding/json"
)

// Decode unmarshals the configuration data into the target struct.
// It serializes the internal data map to JSON and then unmarshals it.
func (c *Config) Decode(target any) error {
	data, err := json.Marshal(c.data)
	if err != nil {
		return WrapDecodeError(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return WrapDecodeError(err)
	}
	return nil
}
