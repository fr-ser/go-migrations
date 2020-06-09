package migrate

import (
	"github.com/urfave/cli/v2"
)

// MigrateCommands perform migration actions
var MigrateCommands = &cli.Command{
	Name:  "migrate",
	Usage: "perform migration actions",
	Subcommands: []*cli.Command{
		migrateUpCommand,
	},
}
