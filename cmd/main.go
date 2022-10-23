package main

import (
	"gochat/cmd/options"
	"log"
	"os"
)

const (
	name  = "gochat"
	usage = "simple and fast im service"
)

func main() {
	// bootstrap with command line
	cmd := options.NewCommand(name, usage)

	if err := cmd.Run(os.Args); err != nil {
		log.Panic(err)
	}
}
