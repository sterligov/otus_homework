package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadbarLine(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		buf := &bytes.Buffer{}
		const (
			dataSize = 450
			lineSize = 5
		)

		loadbar := NewLoadbarLine(dataSize, lineSize, buf)

		tests := []struct {
			progress int
			bar      string
		}{
			{
				progress: 25,
				bar:      "\n\033[F  5%  [     ]\n",
			},
			{
				progress: 175,
				bar:      "\n\033[F  5%  [     ]\n\033[F 44%  [     ]\n\033[F 44%  [==   ]\n",
			},
			{
				progress: 250,
				bar:      "\n\033[F  5%  [     ]\n\033[F 44%  [     ]\n\033[F 44%  [==   ]\n\033[F100%  [==   ]\n\033[F100%  [=====]\n",
			},
		}

		for _, tst := range tests {
			tst := tst
			testName := fmt.Sprintf("progress %d", tst.progress)

			t.Run(testName, func(t *testing.T) {
				loadbar.Update(tst.progress)

				require.Equal(t, tst.bar, buf.String())
			})
		}
	})
}

func TestBarWriter(t *testing.T) {
	t.Run("write data", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bar := NewBarWriter(buf, &loadbarMock{})

		dataToWrite := "abc"
		n, err := bar.Write([]byte(dataToWrite))

		require.NoError(t, err)
		require.Equal(t, len(dataToWrite), n)
		require.Equal(t, dataToWrite, buf.String())
	})

	t.Run("write empty data", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bar := NewBarWriter(buf, &loadbarMock{})

		dataToWrite := ""
		n, err := bar.Write([]byte(dataToWrite))

		require.NoError(t, err)
		require.Equal(t, len(dataToWrite), n)
		require.Equal(t, dataToWrite, buf.String())
	})
}

type loadbarMock struct {
	progress int
}

func (l *loadbarMock) Update(int) {}
