package storage

import "errors"

var (
	insertErr   = errors.New("insert query error")
	notFoundErr = errors.New("no user found")
)
