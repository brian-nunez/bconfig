package main

import (
	"fmt"
	"log"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/file"
)

func file_example() {
	cfg, err := bconfig.New(file.Source("config.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	addr := cfg.String("server.addr")
	fmt.Printf("file: server.addr = %s\n", addr)
}
