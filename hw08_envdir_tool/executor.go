package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

const ErrorExitCode = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, stdin io.Reader, stdout, stderr io.Writer) (returnCode int) {
	for k, v := range env {
		var err error
		if v.NeedRemove {
			err = os.Unsetenv(k)
		} else {
			err = os.Setenv(k, v.Value)
		}

		if err != nil {
			returnCode = ErrorExitCode
			return
		}
	}

	c := cmd[0]
	params := cmd[1:]
	command := exec.Command(c, params...)

	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	command.Env = os.Environ()

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return ErrorExitCode
	}

	return command.ProcessState.ExitCode()
}
