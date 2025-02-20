package storage

import "errors"

var (
	ErrUrlNotFound      = errors.New("url not found")
	ErrAliasNotFound    = errors.New("alias not found")
	ErrUrlNotExists     = errors.New("url not exist")
	ErrUrlAlreadyExists = errors.New("url is already exists")
)
