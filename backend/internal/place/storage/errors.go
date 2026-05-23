package storage

import "errors"

var (
	ErrPlaceNotFound            = errors.New("place not found")
	ErrPlaceAlreadyExists       = errors.New("place already exists")
	ErrCategoryNotFound         = errors.New("category not found")
	ErrPlaceReportAlreadyExists = errors.New("place report already exists")
)
