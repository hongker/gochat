package options

import (
	"github.com/urfave/cli/v2"
	"gochat/internal"
)

// NewCommand returns a new instance of *cli.App
func NewCommand(name, usage string) *cli.App {
	opts := NewServerRunOptions()

	app := &cli.App{
		Name:    name,
		Version: internal.Version,
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
func run(opts *ServerRunOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	config := internal.DefaultConfig()
	opts.applyTo(config)

	return config.New().Run()
}
