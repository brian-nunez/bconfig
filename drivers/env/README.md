# bconfig/drivers/env

`env` is a configuration source driver for the `bconfig` package. It reads environment variables, filters them by prefixes, and parses them as flat keys.

## Installation

```bash
go get github.com/brian-nunez/bconfig/drivers/env
```

## Usage

```go
package main

import (
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/env"
)

func main() {
	// set the environment variable for testing `BAPP_SERVER_ADDR=localhost:8080`
	cfg, err := bconfig.New(
		env.Source("BAPP_"),
		// env.Source(), // to read all environment variables without filtering by prefix
	)
	if err != nil {
		log.Fatal(err)
	}

	addr := cfg.String("bapp_server_addr")
	log.Printf("Server address: %s", addr)
}
```
