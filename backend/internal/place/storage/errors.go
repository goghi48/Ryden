package storage

import "errors"

var (
	ErrPlaceNotFound      = errors.New("place not found")
	ErrPlaceAlreadyExists = errors.New("place already exists")
)
