package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		args []string
		name string
		err  error
		code int
	}{
		{
			args: []string{"command", "arg1=val1", "arg2=val2"},
			name: "ok",
			err:  nil,
			code: ExitOK,
		},
		{
			args: []string{"command", "arg=val"},
			name: "exit error",
			err:  &exec.ExitError{},
			code: -1,
		},
		{
			args: []string{"command"},
			name: "not exit error",
			err:  fmt.Errorf("unknown error"),
			code: ExitFail,
		},
		{
			args: nil,
			name: "empty args",
			err:  nil,
			code: ExitFail,
		},
	}

	env := Environment{
		"FOO": "VAL_FOO",
		"BOO": "VAL_BOO",
		"ZOO": "",
	}

	for _, tst := range tests {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				cmd := &exec.Cmd{
					Path: name,
					Args: append([]string{name}, arg...),
				}

				monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Start", func(_ *exec.Cmd) error {
					require.Equal(t, tst.args[0], cmd.Path)
					require.ElementsMatch(t, tst.args, cmd.Args)
					require.ElementsMatch(t, env.ToArray(), cmd.Env)
					require.Equal(t, cmd.Stdin, os.Stdin)
					require.Equal(t, cmd.Stdout, os.Stdout)
					require.Equal(t, cmd.Stderr, os.Stderr)

					return nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Wait", func(_ *exec.Cmd) error {
					return tst.err
				})

				return cmd
			})

			returnCode := RunCmd(tst.args, env)

			require.Equal(t, tst.code, returnCode)
		})
	}

	t.Run("not_exist_command", func(t *testing.T) {
		monkey.Unpatch(exec.Command)
		returnCode := RunCmd([]string{"not_exist_command"}, nil)

		require.Equal(t, ExitFail, returnCode)
	})
}
