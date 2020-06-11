package migrate

import (
	"fmt"
	"strconv"
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

// migrateUpCommand executes a (up) migrations
var migrateUpCommand = &cli.Command{
	Name:   "up",
	Usage:  "executes a (up) migrations",
	Flags:  flags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {

		if err := checkFlags(c); err != nil {
			return err
		}

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
			if err := db.EnsureConsistentMigrations(); err != nil {
				return err
			}

			upCount := c.Uint("count")
			if upCount == 0 && !c.Bool("all") {
				upCount = 1
			}
			if err := db.ApplyUpMigrationsWithCount(upCount, c.Bool("all")); err != nil {
				return err
			}
		}
		log.Info("Up migration completed")
		return nil
	},
}

func checkFlags(c *cli.Context) error {
	upParamsCount := 0
	if c.String("count") != "" {
		upParamsCount = upParamsCount + 1
	}
	if c.String("only") != "" {
		upParamsCount = upParamsCount + 1
	}
	if c.Bool("all") {
		upParamsCount = upParamsCount + 1
	}

	if upParamsCount > 1 {
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
