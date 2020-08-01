package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"

	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"go-migrations/commands/bootstrap"
	"go-migrations/commands/createseed"
	"go-migrations/commands/migrate"
	"go-migrations/commands/start"
	"go-migrations/utils"
)

func errExitHandler(c *cli.Context, err error) {
	log.Fatal(err)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initLogger()

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.ExitErrHandler = errExitHandler
	app.Commands = []*cli.Command{
		start.StartCommand,
		bootstrap.BootstrapCommand,
		migrate.MigrateCommands,
		createseed.CreateSeedCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		// this should not be called, as we have an exiting error handler
		errExitHandler(nil, err)
	}
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		PadLevelText:     true,
	})

	switch logLevel := utils.GetEnvDefault("LOG_LEVEL", "INFO"); logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		formatter := &log.TextFormatter{
			FullTimestamp: true,
			PadLevelText:  true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf(" %s:%d", filename, f.Line)
			},
		}
		log.SetFormatter(formatter)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
