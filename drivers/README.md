# bconfig Configuration Drivers

This directory contains the configuration source drivers for `bconfig` (a lightweight, composable configuration loader for Go).

Drivers are separate Go packages that implement the [bconfig.Source](../source.go) interface.

---

## Existing Drivers

Each driver is hosted in its own package/sub-module. Click the links below to view their configurations and implementations:

1. **[env](./env)**
   - **Type**: Environment variable loader.
   - **Use Case**: Production environment settings, credentials, port configuration, or 12-factor application flags.
   - **Configuration**: Accepts list of string prefixes to filter variables.
   - **Features**: Automatic lowercase key mapping, filtering by prefixes, and basic type parsing (bool, int, float).

2. **[file](./file)**
   - **Type**: JSON/YAML configuration file loader.
   - **Use Case**: Default configuration files, environment-specific configs, or complex nested configuration structures.
   - **Configuration**: File path (supports `.json`, `.yaml`, and `.yml`).
   - **Features**: Automatic file reading, file type detection via extension, and nesting mapping with dot-path support.

---

## Composing Sources

Unlike traditional registries, `bconfig` does not use dynamic package initialization (`init()` registration). Instead, sources are passed directly as arguments to [bconfig.New](../config.go) or [bconfig.Load](../config.go). 

Sources are evaluated in order from left to right. When keys overlap, values from later sources override values from earlier sources:

```go
package main

import (
	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/env"
	"github.com/brian-nunez/bconfig/drivers/file"
)

func main() {
	cfg, err := bconfig.New(
		file.Source("config.yaml"), // Loaded first (defaults)
		env.Source("BAPP_"),        // Loaded second (overrides config.yaml)
	)
	// ...
}
```

---

## How to Implement a Custom Driver

To add a new configuration source to `bconfig`, follow these steps:

### 1. Implement the `bconfig.Source` Interface
Create a type that implements the [bconfig.Source](../source.go) interface:

```go
package custom

import (
	"context"
	"github.com/brian-nunez/bconfig"
)

type customSource struct {
	// Custom connection clients, paths, or keys
}

func (c *customSource) Name() string {
	return "custom-source-name"
}

func (c *customSource) Load(ctx context.Context) (map[string]any, error) {
	// 1. Fetch config data from your backend/service.
	// 2. Parse data into a map[string]any.
	// 3. Return the map.
	configMap := make(map[string]any)
	return configMap, nil
}
```

### 2. Provide a Helper Constructor
Expose a helper function to easily instantiate the source in application code:

```go
func Source() bconfig.Source {
	return &customSource{}
}
```

---

## Shared Error Handling

When returning errors inside the `Load` method, wrap any transport or parse errors using the utilities in [errors.go](../errors.go) to ensure standard formatting:

- If a driver fails to load, wrap it during loading (handled by the core package via `bconfig.WrapSourceError`).
