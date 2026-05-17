package transport

import (
	"context"

	placesv1 "github.com/goghi48/ryden/gen/go/ryden/places/v1"
	"github.com/goghi48/ryden/internal/place/domain"
	"github.com/goghi48/ryden/internal/place/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	placesv1.UnimplementedPlaceServiceServer

	placeService *service.PlaceService
}

func NewHandler(placeService *service.PlaceService) *Handler {
	return &Handler{
		placeService: placeService,
	}
}

func (h *Handler) CreatePlace(
	ctx context.Context,
	req *placesv1.CreatePlaceRequest,
) (*placesv1.CreatePlaceResponse, error) {
	input := service.CreatePlaceInput{
		Title:           req.GetTitle(),
		Description:     req.GetDescription(),
		Address:         req.GetAddress(),
		City:            req.GetCity(),
		Latitude:        req.GetLatitude(),
		Longitude:       req.GetLongitude(),
		CreatedByUserID: req.GetCreatedByUserId(),
	}

	place, err := h.placeService.CreatePlace(input)

	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return &placesv1.CreatePlaceResponse{
		Place: placeToProto(place),
	}, nil
}

func placeToProto(place domain.Place) *placesv1.Place {
	return &placesv1.Place{
		Id:              place.ID,
		Title:           place.Title,
		Description:     place.Description,
		Address:         place.Address,
		City:            place.City,
		Latitude:        place.Latitude,
		Longitude:       place.Longitude,
		CreatedByUserId: place.CreatedByUserID,
		Status:          placeStatusToProto(place.Status),
		CreatedAt:       timestamppb.New(place.CreatedAt),
		UpdatedAt:       timestamppb.New(place.UpdatedAt),
	}
}

func placeStatusToProto(status domain.PlaceStatus) placesv1.PlaceStatus {
	switch status {
	case domain.StatusPendingReview:
		return placesv1.PlaceStatus_PLACE_STATUS_PENDING_REVIEW
	case domain.StatusApproved:
		return placesv1.PlaceStatus_PLACE_STATUS_APPROVED
	case domain.StatusRejected:
		return placesv1.PlaceStatus_PLACE_STATUS_REJECTED
	case domain.StatusArchived:
		return placesv1.PlaceStatus_PLACE_STATUS_ARCHIVED
	default:
		return placesv1.PlaceStatus_PLACE_STATUS_UNSPECIFIED
	}
}

func (h *Handler) GetPlace(
	ctx context.Context,
	req *placesv1.GetPlaceRequest,
) (*placesv1.GetPlaceResponse, error) {
	place, err := h.placeService.GetPlace(req.GetId())

	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return &placesv1.GetPlaceResponse{
		Place: placeToProto(place),
	}, nil
}

func (h *Handler) ListPlaces(
	ctx context.Context,
	req *placesv1.ListPlacesRequest,
) (*placesv1.ListPlacesResponse, error) {
	places, err := h.placeService.ListPlaces(req.GetCity(), int(req.GetLimit()))
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	protoPlaces := make([]*placesv1.Place, len(places))

	for i, place := range places {
		protoPlaces[i] = placeToProto(place)
	}

	return &placesv1.ListPlacesResponse{
		Places: protoPlaces,
	}, nil
}
