package migrate

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

func TestMain(m *testing.M) {
	app.Commands = []*cli.Command{
		MigrateCommands,
	}
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
