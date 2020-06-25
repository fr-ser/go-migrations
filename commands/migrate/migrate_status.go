package migrate

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands"
)

var statusFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "migrations-path", Aliases: []string{"p"}, Value: "./migrations/zlab",
		Usage: "(relative) path to the folder containing the database migrations",
	},
	&cli.StringFlag{
		Name: "environment", Aliases: []string{"e"}, Value: "development",
		Usage: "Name of the environment and the corresponding configuration",
	},
}

// migrateStatusCommand shows the status of applied and unapplied migrations
var migrateStatusCommand = &cli.Command{
	Name:   "status",
	Usage:  "shows the status of applied and unapplied migrations",
	Flags:  statusFlags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {

		db, err := mockableLoadDB(c.String("migrations-path"), c.String("environment"))
		if err != nil {
			return err
		}

		if err := db.WaitForStart(100*time.Millisecond, 1); err != nil {
			return err
		}

		created, err := db.EnsureMigrationsChangelog()
		if created {
			log.Warning("Created changelog table")
		}
		err = db.PrintMigrationStatus()
		if err != nil {
			return err
		}

		return nil
	},
}
