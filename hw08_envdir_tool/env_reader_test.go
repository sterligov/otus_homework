package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("correct reading", func(t *testing.T) {
		expected := Environment{
			"BAR":   "bar",
			"FOO":   "   foo\nwith new line",
			"HELLO": `"hello"`,
		}
		actual, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("not existing dir", func(t *testing.T) {
		_, err := ReadDir("not_existing_dir")

		require.Error(t, err)
	})

	t.Run("empty dir", func(t *testing.T) {
		dir := "testdata/empty"

		err := os.Mkdir(dir, os.FileMode(0755))
		require.Nil(t, err)
		defer func() {
			err := os.Remove(dir)
			require.NoError(t, err)
		}()

		env, err := ReadDir(dir)

		require.NoError(t, err)
		require.Equal(t, Environment{}, env)
	})

	t.Run("forbidden filename", func(t *testing.T) {
		dir := "testdata/forbidden_env"

		err := os.Mkdir(dir, os.FileMode(0755))
		require.NoError(t, err)
		defer func() {
			err := os.RemoveAll(dir)
			require.NoError(t, err)
		}()

		f, err := os.Create(path.Join(dir, "NAME=FORBIDDEN"))
		require.Nil(t, err)
		defer func() {
			err := f.Close()
			require.NoError(t, err)
		}()

		env, err := ReadDir(dir)

		require.Zero(t, env)
		require.Equal(t, ErrForbiddenFilename, err)
	})
}
