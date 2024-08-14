package storage

import "errors"

var (
	InsertErr   = errors.New("insert query error")
	NotFoundErr = errors.New("no user found")
)
