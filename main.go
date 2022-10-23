package main

import (
	"github.com/ebar-go/ego/utils/runtime"
	"gochat/internal"
	"log"
	"os"
)

const (
	name  = "gochat"
	usage = "simple and fast im service"
)

func main() {
	// bootstrap with command line
	cmd := internal.NewCommand(name, usage)

	// run the command with os.Args.
	runtime.HandleError(cmd.Run(os.Args), func(err error) {
		log.Panic(err)
	})
}
