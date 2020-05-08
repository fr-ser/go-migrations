package start

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/fr-ser/go-migrations/commands"
)

// StartCommand starts a local development database based on a docker-compose file
var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "starts a local development database based on a docker-compose file",
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {
		fmt.Print("Hello")
		return nil
	},
}
