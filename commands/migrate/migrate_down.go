package migrate

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands"
	"go-migrations/internal/direction"
)

var downFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "count", Aliases: []string{"c"},
		Usage: "number of migrations to apply (default action is to apply one)",
	},
	&cli.BoolFlag{
		Name: "all", Aliases: []string{"A"},
		Usage: "apply all outstanding up migrations",
	},
	&cli.StringFlag{
		Name: "only", Aliases: []string{"o"},
		Usage: "apply only one migration containing this string",
	},
	&cli.StringFlag{
		Name: "migrations-path", Aliases: []string{"p"}, Value: "./migrations/zlab",
		Usage: "(relative) path to the folder containing the database migrations",
	},
	&cli.StringFlag{
		Name: "environment", Aliases: []string{"e"}, Value: "development",
		Usage: "Name of the environment and the corresponding configuration",
	},
}

// migrateDownCommand executes down migrations
var migrateDownCommand = &cli.Command{
	Name:   "down",
	Usage:  "executes down migrations",
	Flags:  downFlags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {

		if err := checkFlags(c); err != nil {
			return err
		}

		db, err := mockableLoadDB(c.String("migrations-path"), c.String("environment"))
		if err != nil {
			return err
		}

		if err := db.WaitForStart(100*time.Millisecond, 1); err != nil {
			return err
		}
		log.Info("Connected to database")

		created, err := db.EnsureMigrationsChangelog()
		if created {
			log.Warning("Created changelog table")
		}
		if err != nil {
			return err
		}

		if c.String("only") != "" {
			if err := db.ApplySpecificMigration(c.String("only"), direction.Down); err != nil {
				return err
			}
		} else {
			if err := db.EnsureConsistentMigrations(); err != nil {
				return err
			}

			upCount := c.Uint("count")
			if upCount == 0 && !c.Bool("all") {
				upCount = 1
			}
			err = db.ApplyMigrationsWithCount(upCount, c.Bool("all"), direction.Down)
			if err != nil {
				return err
			}
		}
		log.Info("Down migration completed")
		return nil
	},
}
