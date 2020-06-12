package start

import (
	"fmt"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands"
	"go-migrations/database/driver"
	"go-migrations/utils"
)

// variables to allow mocking for tests
var (
	runWithOutput  = utils.RunWithOutput
	mockableLoadDB = driver.LoadDB
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name: "dc-file", Aliases: []string{"d"}, Value: "docker-compose.yaml",
		Usage: "Path to docker compose file",
	},
	&cli.StringFlag{
		Name: "service", Aliases: []string{"s"}, Value: "database",
		Usage: "service name (in the docker-compose file) of the database",
	},
	&cli.BoolFlag{
		Name: "restart", Aliases: []string{"r"},
		Usage: "stop the docker-compose database service before starting",
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

// StartCommand starts a local development database based on a docker-compose file
var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "starts a local development database based on a docker-compose file",
	Flags:  flags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {
		var err error

		if c.Bool("restart") {
			err = stopDb(c.String("dc-file"), c.String("service"))
			if err != nil {
				return fmt.Errorf("Could not stop database - Err: %v", err)
			}
		}

		err = startDb(c.String("dc-file"), c.String("service"))
		if err != nil {
			return fmt.Errorf("Could not start database - Err: %v", err)
		}

		db, err := mockableLoadDB(c.String("migrations-path"), c.String("environment"))
		if err != nil {
			return err
		}

		if err := db.WaitForStart(1*time.Second, 10); err != nil {
			return err
		}
		log.Info("Connected to database")

		if _, err := db.EnsureMigrationsChangelog(); err != nil {
			return err
		}

		if err := db.Bootstrap(); err != nil {
			return err
		}
		log.Info("Applied bootstrap migration")

		if err := db.ApplyAllUpMigrations(); err != nil {
			return err
		}
		log.Info("Applied all migrations")

		return err
	},
}

func startDb(dcFile, service string) error {
	cmd := exec.Command("docker-compose", "--file", dcFile, "up", "--detach", service)

	_, stderr, err := runWithOutput(cmd)
	if err != nil {
		log.Error(stderr)
		return err
	}

	log.Info("Started database")

	return nil
}

func stopDb(dcFile, service string) error {
	cmd := exec.Command("docker-compose", "--file", dcFile, "rm", "--force", "--stop", service)

	_, stderr, err := runWithOutput(cmd)
	if err != nil {
		log.Error(stderr)
		return err
	}

	log.Debug("Stopped database")

	return nil
}
