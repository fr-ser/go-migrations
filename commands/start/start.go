package start

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fr-ser/go-migrations/commands"
	"github.com/fr-ser/go-migrations/utils"
)

// function variables to allow mocking for tests
var (
	runWithOutput = utils.RunWithOutput
	checkError    = utils.CheckError
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name: "dc-file", Aliases: []string{"d"}, Value: "",
		Usage: "Path to alternate docker compose file",
	},
	&cli.StringFlag{
		Name: "path", Aliases: []string{"p"}, Value: "./migrations",
		Usage: "(relative) path to the folder containing the database migrations",
	},
	&cli.StringFlag{
		Name: "db", Value: "zlab",
		Usage: "name of database migration folder",
	},
}

// StartCommand starts a local development database based on a docker-compose file
var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "starts a local development database based on a docker-compose file",
	Flags:  flags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) error {
		log.Infof("dc-file: %s", c.String("dc-file"))
		log.Infof("path: %s", c.String("path"))
		log.Infof("db: %s", c.String("db"))

		var cmd *exec.Cmd
		if c.String("dc-file") == "" {
			cmd = exec.Command("docker-compose", "up")
		} else {
			cmd = exec.Command("docker-compose", "--file", c.String("dc-file"), "up")
		}

		_, stderr, err := runWithOutput(cmd)
		if err != nil {
			log.Error(stderr)
			checkError("Could not run docker-compose up", err)
			return err
		}

		log.Info("Started docker-compose")

		return nil
	},
}
