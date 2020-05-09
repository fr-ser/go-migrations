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

	"github.com/fr-ser/go-migrations/commands/start"
	"github.com/fr-ser/go-migrations/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	initLogger()

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		start.StartCommand,
	}

	err := app.Run(os.Args)
	utils.CheckError("Could not run the CLI command (top level)", err)
}

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})

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
