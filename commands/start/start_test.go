package start

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

func strSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx++ {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true
}

func TestMain(m *testing.M) {
	app.Commands = []*cli.Command{
		StartCommand,
	}
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
