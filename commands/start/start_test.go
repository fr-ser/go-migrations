package start

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/urfave/cli/v2"
)

type fakeRunWithSpy struct {
	Cmds    []*exec.Cmd
	LastCmd *exec.Cmd
}

var (
	successStdout = "stdout_success"
	successStderr = "stderr_success"
	failStdout    = "stdout_fail"
	failStderr    = "stderr_fail"
)

func (f *fakeRunWithSpy) runWithOutputSuccess(cmd *exec.Cmd) (stdout, stderr string, err error) {
	f.Cmds = append(f.Cmds, cmd)
	f.LastCmd = cmd
	return successStdout, successStderr, nil
}

func (f *fakeRunWithSpy) runWithOutputFail(cmd *exec.Cmd) (stdout, stderr string, err error) {
	f.Cmds = append(f.Cmds, cmd)
	f.LastCmd = cmd
	return failStdout, failStderr, errors.New("Sadly it failed")
}

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

// TODO: Test functionality to run migrations

// TODO: Test path and db flags
// &cli.StringFlag{
// 	Name: "path", Aliases: []string{"p"}, Value: "./migrations",
// 	Usage: "(relative) path to the folder containing the database migrations",
// },
// &cli.StringFlag{
// 	Name: "db", Value: "zlab",
// 	Usage: "name of database migration folder",
// },

func TestStartWithDefaults(t *testing.T) {
	args := []string{"sth.exe", "start"}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputSuccess

	hook := test.NewGlobal()
	defer hook.Reset()

	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"docker-compose", "up", "--detach", "database"}
	if !strSliceEqual(fakeRun.LastCmd.Args, expected) {
		t.Errorf("Expected to run command '%v', but got %s", expected, fakeRun.LastCmd.Args)
	}
	for _, element := range hook.AllEntries() {
		if element.Message == successStderr {
			t.Error("Expected start to NOT log the stderr of the program for a success case")
		}
	}
}

func TestStartAlternateComposeFile(t *testing.T) {
	args := []string{
		"sth.exe", "start",
		"--dc-file", "./docker-compose/non_standard.yaml",
		"--service", "db",
	}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputSuccess

	err := app.Run(args)

	if err != nil {
		t.Errorf("Error running command - %s", err)
	}

	// no restart
	if len(fakeRun.Cmds) != 1 {
		t.Errorf("Expected 1 commands to run but got %d", len(fakeRun.Cmds))
	}

	expected := []string{
		"docker-compose", "--file", "./docker-compose/non_standard.yaml",
		"up", "--detach", "db",
	}
	if !strSliceEqual(fakeRun.LastCmd.Args, expected) {
		t.Errorf("Expected to run command '%v', but got %s", expected, fakeRun.LastCmd.Args)
	}
}

func TestStartWithRestart(t *testing.T) {
	args := []string{"sth.exe", "start", "--restart", "-s", "db"}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputSuccess

	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	if len(fakeRun.Cmds) != 2 {
		t.Errorf("Expected 2 commands to run but got %d", len(fakeRun.Cmds))
	}

	expectedStop := []string{"docker-compose", "rm", "--force", "--stop", "db"}
	if !strSliceEqual(fakeRun.Cmds[0].Args, expectedStop) {
		t.Errorf("Expected to run command '%v', but got %s", expectedStop, fakeRun.Cmds[0].Args)
	}
	expectedStart := []string{"docker-compose", "up", "--detach", "db"}
	if !strSliceEqual(fakeRun.Cmds[1].Args, expectedStart) {
		t.Errorf("Expected to run command '%v', but got %s", expectedStart, fakeRun.Cmds[1].Args)
	}
}

func TestStartFailure(t *testing.T) {
	args := []string{"sth.exe", "start"}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputFail

	hook := test.NewGlobal()
	defer hook.Reset()

	err := app.Run(args)

	if err == nil {
		t.Errorf("The command did not return an error")
	}

	stdoutPrinted := false
	for _, element := range hook.AllEntries() {
		if element.Message == failStderr {
			stdoutPrinted = true
			break
		}
	}
	if !stdoutPrinted {
		t.Error("Expected start to log the stderr of the program for a failure case")
	}
}
