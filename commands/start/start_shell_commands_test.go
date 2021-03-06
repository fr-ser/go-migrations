package start

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"go-migrations/internal"
)

var (
	successStdout = "stdout_success"
	successStderr = "stderr_success"
	failStdout    = "stdout_fail"
	failStderr    = "stderr_fail"
)

type fakeRunWithSpy struct {
	Cmds    []*exec.Cmd
	LastCmd *exec.Cmd
}

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

func TestStartWithDefaults(t *testing.T) {
	args := []string{"sth.exe", "start"}
	var fakeRun fakeRunWithSpy
	mockableRunWithOutput = fakeRun.runWithOutputSuccess

	hook := test.NewGlobal()
	defer hook.Reset()

	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{
		"docker-compose", "--file", "docker-compose.yaml",
		"up", "--detach", "database",
	}
	if !internal.StrSliceEqual(fakeRun.LastCmd.Args, expected) {
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
	mockableRunWithOutput = fakeRun.runWithOutputSuccess

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
	if !internal.StrSliceEqual(fakeRun.LastCmd.Args, expected) {
		t.Errorf("Expected to run command '%v', but got %s", expected, fakeRun.LastCmd.Args)
	}
}

func TestStartWithRestart(t *testing.T) {
	args := []string{"sth.exe", "start", "--restart", "-s", "db"}
	var fakeRun fakeRunWithSpy
	mockableRunWithOutput = fakeRun.runWithOutputSuccess

	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	if len(fakeRun.Cmds) != 2 {
		t.Errorf("Expected 2 commands to run but got %d", len(fakeRun.Cmds))
	}

	expectedStop := []string{
		"docker-compose", "--file", "docker-compose.yaml",
		"rm", "--force", "--stop", "db",
	}
	if !internal.StrSliceEqual(fakeRun.Cmds[0].Args, expectedStop) {
		t.Errorf("Expected to run command '%v', but got %s", expectedStop, fakeRun.Cmds[0].Args)
	}
	expectedStart := []string{
		"docker-compose", "--file", "docker-compose.yaml",
		"up", "--detach", "db",
	}
	if !internal.StrSliceEqual(fakeRun.Cmds[1].Args, expectedStart) {
		t.Errorf("Expected to run command '%v', but got %s", expectedStart, fakeRun.Cmds[1].Args)
	}
}

func TestStartFailure(t *testing.T) {
	args := []string{"sth.exe", "start"}
	var fakeRun fakeRunWithSpy
	mockableRunWithOutput = fakeRun.runWithOutputFail

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
