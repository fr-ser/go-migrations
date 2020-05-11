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
		Name: "service", Aliases: []string{"s"}, Value: "database",
		Usage: "service name (in the docker-compose file) of the database",
	},
	&cli.BoolFlag{
		Name: "restart", Aliases: []string{"r"},
		Usage: "stop the docker-compose database service before starting",
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
			checkError("Could not stop database", err)
		}

		err = startDb(c.String("dc-file"), c.String("service"))
		checkError("Could not start database", err)
		return err
	},
}

func startDb(dcFile, service string) error {
	args := []string{}
	if dcFile != "" {
		args = append(args, "--file", dcFile)
	}
	args = append(args, "up", "--detach", service)

	_, stderr, err := runWithOutput(exec.Command("docker-compose", args...))
	if err != nil {
		log.Error(stderr)
		checkError("Could not run docker-compose up", err)
		return err
	}

	log.Info("Started database")

	return nil
}

func stopDb(dcFile, service string) error {
	args := []string{}
	if dcFile != "" {
		args = append(args, "--file", dcFile)
	}
	args = append(args, "rm", "--force", "--stop", service)

	_, stderr, err := runWithOutput(exec.Command("docker-compose", args...))
	if err != nil {
		log.Error(stderr)
		checkError("Could not run docker-compose up", err)
		return err
	}

	log.Debug("Stopped database")

	return nil
}
