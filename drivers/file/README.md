# bconfig/drivers/file

`file` is a configuration source driver for the `bconfig` package. It reads JSON, YAML, and YML files and parses them into hierarchical maps.

## Installation

```bash
go get github.com/brian-nunez/bconfig/drivers/file
```

## Usage

```go
package main

import (
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/file"
)

func main() {
	// Sample yaml file
	// server:
	//   addr: localhost:8080
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
