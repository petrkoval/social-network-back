package storage

import "errors"

var (
	NotFoundUserErr  = errors.New("no user found")
	NotFoundTokenErr = errors.New("no token found")

	NotFoundChannelErr = errors.New("no channel found")
)
