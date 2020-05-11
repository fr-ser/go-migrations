package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

// RunWithOutput executes a command and stores the stdout and stderr.
// The returned Error can be an os/exec.ExitError
func RunWithOutput(cmd *exec.Cmd) (stdout, stderr string, err error) {

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("Could not create StderrPipe - %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("Could not create StdoutPipe - %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("Could not start cmd - %v", err)
	}

	stdoutBytes, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return "", "", fmt.Errorf("Could not get stdout of cmd - %v", err)
	}
	stdout = string(stdoutBytes)

	stderrBytes, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		return "", "", fmt.Errorf("Could not get stderr of cmd - %v", err)
	}
	stderr = string(stderrBytes)

	err = cmd.Wait()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return stdout, stderr, err
	} else if err != nil {
		return stdout, stderr, fmt.Errorf("Error waiting for cmd to finish - %v", err)
	}

	return stdout, stderr, nil
}
