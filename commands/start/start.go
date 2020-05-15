package start

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands"
	"go-migrations/utils"
)

// variables to allow mocking for tests
var (
	runWithOutput = utils.RunWithOutput
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
			if err != nil {
				log.Errorf("Could not stop database - Err: %v", err)
				return err
			}
		}

		err = startDb(c.String("dc-file"), c.String("service"))
		if err != nil {
			log.Errorf("Could not start database - Err: %v", err)
			return err
		}

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
		return err
	}

	log.Debug("Stopped database")

	return nil
}
