package bootstrap

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands"
	"go-migrations/database/driver"
)

var (
	mockableLoadDB = driver.LoadDB
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name: "migrations-path", Aliases: []string{"p"}, Value: "./migrations/zlab",
		Usage: "(relative) path to the folder containing the database migrations",
	},
	&cli.StringFlag{
		Name: "environment", Aliases: []string{"e"}, Value: "development",
		Usage: "Name of the environment and the corresponding configuration",
	},
}

// BootstrapCommand bootstraps an already running (empty) database
var BootstrapCommand = &cli.Command{
	Name:   "bootstrap",
	Usage:  "bootstraps an already running (empty) database",
	Flags:  flags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {
		var err error

		db, err := mockableLoadDB(c.String("migrations-path"), c.String("environment"))
		if err != nil {
			return err
		}

		if err := db.WaitForStart(1*time.Second, 10); err != nil {
			return err
		}
		log.Debug("Connected to database")

		if _, err := db.EnsureMigrationsChangelog(); err != nil {
			return err
		}

		if err := db.Bootstrap(); err != nil {
			return err
		}
		log.Info("Applied bootstrap migration")

		pw := progress.NewWriter()
		pw.SetAutoStop(true)
		pw.SetTrackerPosition(progress.PositionRight)
		pw.SetUpdateFrequency(time.Millisecond * 250)

		go pw.Render()

		if err := db.ApplyAllUpMigrations(pw); err != nil {
			return err
		}
		log.Debug("Applied all migrations")

		if pw.IsRenderInProgress() {
			time.Sleep(time.Millisecond * 251)
		}

		return nil
	},
}
