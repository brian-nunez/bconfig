# bconfig

`bconfig` is a lightweight configuration loading library for go that supports composing and overriding configuration from multiple sources using a clean driver/source pattern.

---

## What `bconfig` Is
* A library to **load, merge, read, and decode** configurations.
* Supports combining configuration from files (YAML, JSON) and environment variables.
* Supports hierarchical structures with dot-path lookups (specifically for file configurations).
* Handles deep merging where later sources override earlier sources.
* Supports decoding configurations into Go structs.

## What `bconfig` Is Not
`bconfig` is strictly a configuration loading library. It does not provide:
* HTTP/admin APIs or routing
* KV database setup or DB storage
* Dependency injection
* Graceful shutdown or service lifecycle management
* Tracing, metrics, or telemetry

---

## Installation

```bash
go get github.com/brian-nunez/bconfig
```

---

## Usage Examples

### Basic File Example

Load configuration from a YAML or JSON file:

```go
package main

import (
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/file"
)

func main() {
	cfg, err := bconfig.New(
		file.Source("config.yaml"),
	)
	if err != nil {
		log.Fatal(err)
	}

	addr := cfg.String("server.addr")
	log.Printf("Server address: %s", addr)
}
```

### Basic Env Example

Load configuration from environment variables. Keys are lowercased and kept as flat keys:

```go
package main

import (
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/env"
)

func main() {
	// If SERVER_PORT=8080, it maps to server_port
	cfg, err := bconfig.New(
		env.Source(),
	)
	if err != nil {
		log.Fatal(err)
	}

	port := cfg.Int("server_port")
	log.Printf("Server port: %d", port)
}
```

### File + Env Example

Chain sources in order. Prefixes are never stripped from key names (preventing data loss), and environment variables are mapped as flat keys:

```go
package main

import (
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/env"
	"github.com/brian-nunez/bconfig/drivers/file"
)

func main() {
	cfg, err := bconfig.New(
		file.Source("config.yaml"),
		env.Source("BAPP_"), // Loads only variables starting with BAPP_
	)
	if err != nil {
		log.Fatal(err)
	}

	// Read environment variable BAPP_SERVER_ADDR
	addr := cfg.String("bapp_server_addr")
	log.Printf("Server address: %s", addr)
}
```

### Decode into Struct Example

```go
package main

import (
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/file"
)

type AppConfig struct {
	Server struct {
		Addr string `json:"addr" yaml:"addr"`
	} `json:"server" yaml:"server"`
}

func main() {
	cfg, err := bconfig.New(
		file.Source("config.yaml"),
	)
	if err != nil {
		log.Fatal(err)
	}

	var appCfg AppConfig
	if err := cfg.Decode(&appCfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("Decoded server address: %s", appCfg.Server.Addr)
}
```

---

## Merge Order

Sources are evaluated sequentially from left to right. When multiple sources define the same key path, later sources override earlier sources.

Example:

```go
cfg, err := bconfig.New(
	file.Source("config.yaml"),
	env.Source("BAPP_"),
)
```

Merge order:
```
config.yaml
  ↳ BAPP_ environment variables
```

If `config.yaml` defines:
```yaml
server:
  addr: ":8080"
  timeout: "5s"
```

And `BAPP_SERVER_ADDR=:9090` is set in the environment:
1. The `file.Source` loads the initial configuration.
2. The `env.Source` loads `BAPP_SERVER_ADDR=:9090` (keeping prefix and mapping flatly).
3. The values are deep-merged, resulting in:
```yaml
server:
  addr: ":8080"
  timeout: "5s"
bapp_server_addr: ":9090"
```

---

## Future Sources (Planned)

The package is designed with a core `Source` interface, allowing new sources to be added without modifying the root package.

### Vault
```go
vault.Source(vault.Options{
	Addr:  "https://vault.local",
	Token: token,
	Path:  "secret/data/my-app",
})
```

### Consul
```go
consul.Source(consul.Options{
	Addr: "localhost:8500",
	Key:  "apps/my-app/config",
})
```

### S3
```go
s3.Source(s3.Options{
	Bucket: "configs",
	Key:    "my-app/config.yaml",
	Region: "us-east-1",
})
```
