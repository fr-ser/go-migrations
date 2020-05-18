package start

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

func TestMain(m *testing.M) {
	app.Commands = []*cli.Command{
		StartCommand,
	}
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
