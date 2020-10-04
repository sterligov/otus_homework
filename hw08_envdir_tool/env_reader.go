package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var ErrForbiddenFilename = errors.New("forbidden filename with =")

type Environment map[string]string

func (e Environment) ToArray() []string {
	env := make([]string, 0, len(e))
	for k, v := range e {
		env = append(env, k+"="+v)
	}

	return env
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(filesInfo))

	for _, fi := range filesInfo {
		if fi.IsDir() {
			continue
		}

		if strings.Contains(fi.Name(), "=") {
			return nil, ErrForbiddenFilename
		}

		if fi.Size() == 0 {
			delete(env, fi.Name())
			continue
		}

		fname := path.Join(dir, fi.Name())
		f, err := os.Open(fname)
		if err != nil {
			return nil, err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Println(err)
			}
		}()

		buf := bufio.NewReader(f)
		line, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		val := bytes.ReplaceAll(line, []byte{0}, []byte{'\n'})
		env[fi.Name()] = strings.TrimRight(string(val), " \t\n")
	}

	return env, nil
}
