package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestCopy(t *testing.T) {
	t.Run("copy", func(t *testing.T) {
		tests := []struct {
			limit  int
			offset int
			err    error
		}{
			{
				limit:  0,
				offset: 0,
				err:    nil,
			},
			{
				limit:  0,
				offset: 15,
				err:    nil,
			},
			{
				limit:  10,
				offset: 50,
				err:    nil,
			},
			{
				limit:  2048,
				offset: 0,
				err:    nil,
			},
			{
				limit:  100,
				offset: 0,
				err:    nil,
			},
			{
				limit:  -1,
				offset: 0,
				err:    nil,
			},
			{
				limit:  0,
				offset: -1,
				err:    ErrOffsetExceedsFileSize,
			},
			{
				limit:  2048,
				offset: 4096,
				err:    ErrOffsetExceedsFileSize,
			},
		}

		const fileSize = 1024

		dstFilename := "dst"
		defer func() {
			err := os.Remove(dstFilename)
			require.NoError(t, err)
		}()

		for _, tst := range tests {
			tst := tst
			testName := fmt.Sprintf("limit:%d offset:%d", tst.limit, tst.offset)

			t.Run(testName, func(t *testing.T) {
				srcFilename, srcData, err := generateRandomFile(fileSize)
				require.NoError(t, err)
				defer func() {
					err := os.Remove(srcFilename)
					require.NoError(t, err)
				}()

				err = Copy(srcFilename, dstFilename, int64(tst.offset), int64(tst.limit))
				require.Equal(t, tst.err, err)
				if err != nil {
					return
				}

				out, err := os.OpenFile(dstFilename, os.O_RDONLY, os.FileMode(0755))
				require.NoError(t, err)
				defer func() {
					err := out.Close()
					require.NoError(t, err)
				}()

				dstData := make([]byte, fileSize)
				_, err = out.Read(dstData)
				require.NoError(t, err)

				srcEnd := fileSize
				if tst.offset+tst.limit < fileSize && tst.limit > 0 {
					srcEnd = tst.limit + tst.offset
				}
				dstEnd := srcEnd - tst.offset

				require.Equal(t, srcData[tst.offset:srcEnd], string(dstData[:dstEnd]))
			})
		}
	})

	t.Run("file not exist", func(t *testing.T) {
		err := Copy("not_exist", "out", 0, 0)
		require.Error(t, err)
	})

	t.Run("not regular file", func(t *testing.T) {
		err := Copy("./", "out", 0, 0)
		require.Error(t, err)
	})
}

func generateRandomFile(size int) (filename string, data string, rerr error) {
	var builder strings.Builder

	for i := 0; i < size; i++ {
		builder.WriteByte('a' + byte(rand.Intn(26)))
	}

	out, rerr := ioutil.TempFile("./", "")
	if rerr != nil {
		return
	}
	defer func() {
		err := out.Close()
		if err != nil {
			rerr = err
		}
	}()

	out.Write([]byte(builder.String()))

	return out.Name(), builder.String(), rerr
}
