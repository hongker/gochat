package options

import (
	"github.com/urfave/cli/v2"
	"gochat/internal"
)

type ServerRunOptions struct {
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{}
}

func (o *ServerRunOptions) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (o *ServerRunOptions) Parse(ctx *cli.Context) {}

func (o *ServerRunOptions) Validate() error {
	return nil
}

func (o *ServerRunOptions) applyTo(config *internal.Config) {}
