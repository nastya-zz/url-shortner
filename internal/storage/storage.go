package storage

import "errors"

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlNotExist = errors.New("url not exist")
)
