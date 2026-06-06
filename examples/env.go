package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/env"
)

func env_example() {
	os.Setenv("BAPP_SERVER_ADDR", ":9090")
	defer os.Unsetenv("BAPP_SERVER_ADDR")

	cfg, err := bconfig.New(env.Source("BAPP_"))
	if err != nil {
		log.Fatal(err)
	}

	addr := cfg.String("bapp_server_addr")
	fmt.Printf("env: bapp_server_addr = %s\n", addr)
}
