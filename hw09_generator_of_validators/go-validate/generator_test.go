package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update validator file")

func TestGenerateValidators(t *testing.T) {
	t.Run("compare with golden file", func(t *testing.T) {
		err := GenerateValidators("./testdata/example.go")
		require.NoError(t, err)

		validatorFilename := "./testdata/example_validation_generated.go"
		actual, err := ioutil.ReadFile(validatorFilename)
		require.NoError(t, err)

		goldenFile := "./testdata/example_validation_golden.go"
		if *update {
			err = ioutil.WriteFile(goldenFile, actual, 0644)
			require.NoError(t, err)
		}

		expected, err := ioutil.ReadFile(goldenFile)
		require.NoError(t, err)

		require.Equal(t, expected, actual)
		require.NoError(t, os.Remove(validatorFilename))
	})

	t.Run("file without validate tag", func(t *testing.T) {
		filename := "./testdata/empty.go"
		err := GenerateValidators(filename)
		require.NoError(t, err)

		_, err = os.Stat("./testdata/empty_example_validation_generated.go")
		require.True(t, os.IsNotExist(err))
	})

	t.Run("not existing file", func(t *testing.T) {
		err := GenerateValidators("not_existing_file")
		require.Error(t, err)
	})

	t.Run("empty params in tag", func(t *testing.T) {
		err := GenerateValidators("./testdata/invalid_params.go")
		require.Error(t, err)
	})
}
