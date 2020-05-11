package start

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/urfave/cli/v2"
)

func mockCheckError(_ string, _ error) {}

type fakeRunWithSpy struct {
	LastCmd *exec.Cmd
}

func (f *fakeRunWithSpy) runWithOutputSuccess(cmd *exec.Cmd) (stdout, stderr string, err error) {
	f.LastCmd = cmd
	return "stdout", "stderr", nil
}

func (f *fakeRunWithSpy) runWithOutputFail(cmd *exec.Cmd) (stdout, stderr string, err error) {
	f.LastCmd = cmd
	return "stdout", "stderr", errors.New("Sadly it failed")
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
	os.Exit(m.Run())
}

// TODO: Add functionality to run migrations

func TestStartDefaultDockerCompose(t *testing.T) {
	args := []string{
		"sth.exe", "start",
		"--path", "/sth/correct",
		"--db", "folder_of_sth",
	}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputSuccess

	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"docker-compose", "up"}
	if !strSliceEqual(fakeRun.LastCmd.Args, expected) {
		t.Errorf("Expected to run command '%v', but got %s", expected, fakeRun.LastCmd.Path)
	}

}

func TestStartAlternateComposeFile(t *testing.T) {
	args := []string{
		"sth.exe", "start",
		"--dc-file", "./docker-compose/non_standard.yaml",
		"--path", "/sth/correct",
		"--db", "folder_of_sth",
	}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputSuccess

	err := app.Run(args)

	if err != nil {
		t.Errorf("Error running command - %s", err)
	}
	expected := []string{"docker-compose", "--file", "./docker-compose/non_standard.yaml", "up"}
	if !strSliceEqual(fakeRun.LastCmd.Args, expected) {
		t.Errorf("Expected to run command '%v', but got %s", expected, fakeRun.LastCmd.Path)
	}
}

func TestStartFailure(t *testing.T) {
	args := []string{"sth.exe", "start"}
	var fakeRun fakeRunWithSpy
	runWithOutput = fakeRun.runWithOutputFail

	originalCheckError := checkError
	defer func() { checkError = originalCheckError }()
	checkError = mockCheckError

	hook := test.NewGlobal()
	defer hook.Reset()

	err := app.Run(args)

	if err == nil {
		t.Errorf("The command did not return an error")
	}

	if hook.LastEntry().Message != "stderr" {
		t.Errorf(
			"Expected start to log the stderr of the program, but got '%v'",
			hook.LastEntry().Message,
		)
	}
}
