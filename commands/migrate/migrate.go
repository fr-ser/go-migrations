package migrate

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"

	"go-migrations/database/driver"
)

// variables to allow mocking for tests
var (
	mockableLoadDB = driver.LoadDB
)

func checkFlags(c *cli.Context) error {
	paramsCount := 0
	if c.String("count") != "" {
		paramsCount = paramsCount + 1
	}
	if c.String("only") != "" {
		paramsCount = paramsCount + 1
	}
	if c.Bool("all") {
		paramsCount = paramsCount + 1
	}

	if paramsCount > 1 {
		return fmt.Errorf(
			"Cannot provide more than one of count (%d), 'only' (%s) and all (%t)",
			c.Uint("count"), c.String("only"), c.Bool("all"),
		)
	}

	if c.String("count") != "" {
		value, err := strconv.ParseUint(c.String("count"), 10, 64)
		if err != nil || value < 1 {
			return fmt.Errorf("Could not format count (%s) to positive integer", c.String("count"))
		}

	}

	return nil
}

// MigrateCommands perform migration actions
var MigrateCommands = &cli.Command{
	Name:  "migrate",
	Usage: "perform migration actions",
	Subcommands: []*cli.Command{
		migrateUpCommand,
		migrateDownCommand,
	},
}
