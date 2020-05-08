package commands

import (
	"errors"

	"github.com/urfave/cli/v2"
)

// NoArguments exits the program if an argument was passed
func NoArguments(c *cli.Context) error {
	if c.NArg() > 0 {
		return errors.New("This command does not support arguments")
	}
	return nil
}
