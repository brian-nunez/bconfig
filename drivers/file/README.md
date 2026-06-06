# bconfig File Driver

`file` is a configuration source driver for the `bconfig` package. It reads, parses, and maps configuration settings from JSON and YAML files.

---

## Installation

```bash
go get github.com/brian-nunez/bconfig/drivers/file
```

## Features

- **Format Autodetection**: Detects file formats automatically via the file extension (`.json`, `.yaml`, and `.yml`).
- **Hierarchical Maps**: Parses nested config files into recursive maps, supporting structured lookup schemes like dot-notation (e.g., `server.addr`).
- **YAML & JSON Support**: Utilizes the robust parser libraries `encoding/json` and `gopkg.in/yaml.v3`.

## Configuration

The driver uses the [file.Config](./config.go) struct:

| Field | Type | Description |
|---|---|---|
| `Path` | `string` | Absolute or relative path to the configuration file (e.g. `"config.yaml"`). |

```go
type Config struct {
	Path string
}
```

## Usage Example

```go
package main

import (
	"context"
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/file"
)

func main() {
	// Loads configuration file from the specified path
	cfg, err := bconfig.New(
		file.Source("config.yaml"),
	)
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Reads nested values using dot-notation path
	addr := cfg.String("server.addr")
	log.Printf("Server Address: %s", addr)
}
```
