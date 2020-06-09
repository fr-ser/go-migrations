package migrate

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands"
	"go-migrations/database/driver"
)

// variables to allow mocking for tests
var (
	dbLoadDb = driver.LoadDb
)

var flags = []cli.Flag{
	&cli.IntFlag{
		Name: "count", Aliases: []string{"c"}, Value: 1,
		Usage: "number of migrations to apply starting from the last applied",
	},
	&cli.BoolFlag{
		Name: "all", Aliases: []string{"A"}, Value: false,
		Usage: "apply all outstanding up migrations (starting from the last applied)",
	},
	&cli.StringFlag{
		Name: "only", Aliases: []string{"o"}, Value: "",
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

// migrateUpCommand executes a (up) migrations
var migrateUpCommand = &cli.Command{
	Name:   "up",
	Usage:  "executes a (up) migrations",
	Flags:  flags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {

		db, err := dbLoadDb(c.String("migrations-path"), c.String("environment"))
		if err != nil {
			return err
		}

		if err := db.WaitForStart(100*time.Millisecond, 1); err != nil {
			return err
		}
		log.Info("Connected to database")

		created, err := db.EnsureMigrationsChangelog()
		if created {
			log.Info("Created changelog table")
		}
		if err != nil {
			return err
		}

		if c.String("only") != "" {
			if err := db.ApplySpecificUpMigration(c.String("only")); err != nil {
				return err
			}
		} else {
			if err := db.ApplyUpMigrationsWithCount(c.Int("count"), c.Bool("all")); err != nil {
				return err
			}
		}
		log.Info("Up migration completed")
		return nil
	},
}
