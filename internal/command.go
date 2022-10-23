package internal

import (
	"github.com/urfave/cli/v2"
	"gochat/internal/options"
)

// NewCommand returns a new instance of *cli.App
func NewCommand(name, usage string) *cli.App {
	opts := options.NewServerRunOptions()

	app := &cli.App{
		Name:    name,
		Version: Version,
		Usage:   usage,
		Flags:   opts.Flags(),
		Action: func(ctx *cli.Context) error {
			opts.Parse(ctx)

			return run(opts)
		},
	}
	return app
}

// run executes command.
func run(opts *options.ServerRunOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	return nil
}
