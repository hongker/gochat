package main

import (
	"embed"
	"gochat/cmd/server/options"
	"gochat/internal/http"
	"log"
	"os"
)

const (
	name  = "gochat"
	usage = "simple and fast im service"
)

var (
	//go:embed app/dist
	Static embed.FS
)

func main() {
	http.Static = Static
	// bootstrap with command line
	cmd := options.NewCommand(name, usage)

	if err := cmd.Run(os.Args); err != nil {
		log.Panic(err)
	}
}
