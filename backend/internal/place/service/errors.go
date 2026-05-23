package service

import "errors"

var (
	ErrTitleRequired           = errors.New("title is required")
	ErrCityRequired            = errors.New("city is required")
	ErrInvalidLatitude         = errors.New("latitude is invalid")
	ErrInvalidLongitude        = errors.New("longitude is invalid")
	ErrInvalidCreatedByUserID  = errors.New("created by user id is invalid")
	ErrInvalidCategoryID       = errors.New("category id is invalid")
	ErrInvalidPlaceID          = errors.New("place id is invalid")
	ErrInvalidReportedByUserID = errors.New("reported by user id is invalid")
	ErrInvalidReportReason     = errors.New("report reason is invalid")
	ErrReportCommentTooLong    = errors.New("comment is too long")
)
