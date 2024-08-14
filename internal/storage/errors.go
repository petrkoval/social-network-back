package storage

import "errors"

var (
	InsertErr        = errors.New("insert query error")
	NotFoundUserErr  = errors.New("no user found")
	NotFoundTokenErr = errors.New("no token found")
)
