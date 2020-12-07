package storage

import "errors"

var (
	ErrNotFound = errors.New("entity not found")
	ErrDateBusy = errors.New("date already busy")
)
