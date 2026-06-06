# bconfig Environment Variable Driver

`env` is a configuration source driver for the `bconfig` package. It loads environment variables, filters them by custom prefixes, and maps them directly as flat keys.

---

## Installation

```bash
go get github.com/brian-nunez/bconfig/drivers/env
```

## Features

- **Optional Prefix Filtering**: Load all environment variables by default or filter for specific variables matching one or more prefixes.
- **Collision Prevention**: Prefix filters are retained on the keys to prevent name collision with general system variables.
- **Case-Insensitive Normalization**: Automatically converts all environment variable keys to lowercase.
- **Smart Type Parsing**: Parses string values into appropriate basic Go types:
  - `"true"` or `"false"` -> `bool`
  - Numeric integer strings -> `int`
  - Float strings -> `float64`
  - Other values -> `string`

## Configuration

The driver uses the [env.Config](./config.go) struct internally:

```go
type Config struct {
	Prefixes []string
}
```

## Usage Example

```go
package main

import (
	"context"
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/env"
)

func main() {
	// Assume BAPP_PORT=8080 is set in the environment.
	// env.Source matches prefix and returns "bapp_port".
	cfg, err := bconfig.New(
		env.Source("BAPP_"),
	)
	if err != nil {
		log.Fatalf("Failed to load environment config: %v", err)
	}

	port := cfg.Int("bapp_port")
	log.Printf("Server port: %d", port)
}
```
