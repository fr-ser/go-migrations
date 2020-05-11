package utils

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	switch os.Getenv("TEST_SUBPROCESS") {
	case "pass":
		os.Stderr.WriteString("stderr")
		os.Stdout.WriteString("stdout")
		os.Exit(0)
	case "fail":
		os.Stderr.WriteString("stderr")
		os.Stdout.WriteString("stdout")
		os.Exit(1)
	case "fail-quiet":
		os.Exit(1)
	default:
		os.Exit(m.Run())
	}
}

func TestRunWithOutputSuccess(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestRunWithOutputSuccess")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=pass")

	stdout, stderr, err := RunWithOutput(cmd)
	if err != nil {
		t.Fatalf("Got error running passing process: %v", err)
	}

	if stdout != "stdout" {
		t.Errorf("Expected stdout of 'stdout' but got: %s", stdout)
	}
	if stderr != "stderr" {
		t.Errorf("Expected stderr of 'stderr' but got: %s", stderr)
	}
}

func TestRunWithOutputFailedProcess(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestRunWithOutputFailedProcess")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=fail")

	stdout, stderr, err := RunWithOutput(cmd)

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("Got error, which is not an ExitError: %v", err)
	}

	if stdout != "stdout" {
		t.Errorf("Expected stdout of 'stdout' but got: %s", stdout)
	}
	if stderr != "stderr" {
		t.Errorf("Expected stderr of 'stderr' but got: %s", stderr)
	}
}
func TestRunWithOutputFailedProcessQuiet(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestRunWithOutputFailedProcessQuiet")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=fail-quiet")

	stdout, stderr, err := RunWithOutput(cmd)

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("Got error, which is not an ExitError: %v", err)
	}

	if stdout != "" {
		t.Errorf("Expected stdout of '' but got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected stderr of '' but got: %s", stderr)
	}
}
func TestRunWithOutputWrongProgram(t *testing.T) {
	cmd := exec.Command("this_thing_should_not_exist.exe")
	stdout, stderr, err := RunWithOutput(cmd)

	if err == nil {
		t.Fatalf("Got no error, but expected executable not found")
	}

	if stdout != "" {
		t.Errorf("Expected stdout of '' but got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected stderr of '' but got: %s", stderr)
	}
}
