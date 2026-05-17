package service

import "errors"

var (
	ErrTitleRequired    = errors.New("title is required")
	ErrCityRequired     = errors.New("city is required")
	ErrInvalidLatitude  = errors.New("latitude is invalid")
	ErrInvalidLongitude = errors.New("longitude is invalid")
)
