package main

import (
	"os"
	"os/exec"
)

const (
	ExitOK = iota
	ExitFail
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(args []string, env Environment) (returnCode int) {
	if len(args) == 0 {
		return ExitFail
	}

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, env.ToArray()...)

	err := cmd.Start()
	if err != nil {
		return ExitFail
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode()
		}

		return ExitFail
	}

	return ExitOK
}
