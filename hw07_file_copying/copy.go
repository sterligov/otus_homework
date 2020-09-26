package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrUnsupportedFile       = errors.New("unsupported file")
)

func Copy(fromPath string, toPath string, offset, limit int64) (rerr error) {
	from, err := os.OpenFile(fromPath, os.O_RDONLY, os.FileMode(0755))
	if err != nil {
		return err
	}
	defer func() {
		err := from.Close()
		if err != nil {
			rerr = err
		}
	}()

	fromStat, err := from.Stat()
	if err != nil {
		return err
	}

	if !fromStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fromStat.Size() || offset < 0 {
		return ErrOffsetExceedsFileSize
	}

	to, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0755))
	if err != nil {
		return err
	}
	defer func() {
		err := to.Close()
		if err != nil {
			rerr = err
		}
	}()

	_, err = from.Seek(offset, 0)
	if err != nil {
		return err
	}

	nByteToCopy := fromStat.Size() - offset
	if limit > 0 && limit < nByteToCopy {
		nByteToCopy = limit
	}

	loadbar := NewLoadbarLine(nByteToCopy, 30 /*bar line size*/, os.Stdout)
	dst := NewBarWriter(to, loadbar)

	_, err = io.CopyN(dst, from, nByteToCopy)
	if err != nil && err != io.EOF {
		return err
	}

	return rerr
}
