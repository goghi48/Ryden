package transport

import (
	"errors"

	"github.com/goghi48/ryden/internal/place/service"
	"github.com/goghi48/ryden/internal/place/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErrorToStatus(err error) error {
	switch {
	case errors.Is(err, service.ErrTitleRequired),
		errors.Is(err, service.ErrCityRequired),
		errors.Is(err, service.ErrInvalidLatitude),
		errors.Is(err, service.ErrInvalidLongitude):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, storage.ErrPlaceNotFound):
		return status.Error(codes.NotFound, err.Error())

	default:
		return status.Error(codes.Internal, "internal error")
	}
}
