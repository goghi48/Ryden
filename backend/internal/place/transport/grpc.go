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

func (h *Handler) CreatePlace(ctx context.Context, req *placesv1.CreatePlaceRequest) (*placesv1.CreatePlaceResponse, error) {
	input := service.CreatePlaceInput{
		Title:           req.GetTitle(),
		Description:     req.GetDescription(),
		Address:         req.GetAddress(),
		City:            req.GetCity(),
		Latitude:        req.GetLatitude(),
		Longitude:       req.GetLongitude(),
		CreatedByUserID: req.GetCreatedByUserId(),
		CategoryIDs:     req.GetCategoryIds(),
	}

	place, err := h.placeService.CreatePlace(ctx, input)

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
		Categories:      categoriesToProto(place.Categories),
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

func (h *Handler) GetPlace(ctx context.Context, req *placesv1.GetPlaceRequest) (*placesv1.GetPlaceResponse, error) {
	place, err := h.placeService.GetPlace(ctx, req.GetId())

	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return &placesv1.GetPlaceResponse{
		Place: placeToProto(place),
	}, nil
}

func (h *Handler) ListPlaces(ctx context.Context, req *placesv1.ListPlacesRequest) (*placesv1.ListPlacesResponse, error) {
	places, err := h.placeService.ListPlaces(ctx, req.GetCity(), int(req.GetLimit()))
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

func (h *Handler) ApprovePlace(ctx context.Context, req *placesv1.ApprovePlaceRequest) (*placesv1.ApprovePlaceResponse, error) {
	place, err := h.placeService.ApprovePlace(ctx, req.GetId())
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return &placesv1.ApprovePlaceResponse{
		Place: placeToProto(place),
	}, nil
}

func (h *Handler) ArchivePlace(ctx context.Context, req *placesv1.ArchivePlaceRequest) (*placesv1.ArchivePlaceResponse, error) {
	place, err := h.placeService.ArchivePlace(ctx, req.GetId())
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return &placesv1.ArchivePlaceResponse{
		Place: placeToProto(place),
	}, nil
}

func (h *Handler) CreatePlaceReport(ctx context.Context, req *placesv1.CreatePlaceReportRequest) (*placesv1.CreatePlaceReportResponse, error) {
	input := service.CreatePlaceReportInput{
		PlaceID:          req.GetPlaceId(),
		ReportedByUserID: req.GetReportedByUserId(),
		Reason:           protoReportReasonToDomain(req.GetReason()),
		Comment:          req.GetComment(),
	}

	report, err := h.placeService.CreatePlaceReport(ctx, input)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return &placesv1.CreatePlaceReportResponse{
		PlaceReport: placeReportToProto(report),
	}, nil
}

func protoReportReasonToDomain(reason placesv1.PlaceReportReason) domain.PlaceReportReason {
	switch reason {
	case placesv1.PlaceReportReason_PLACE_REPORT_REASON_SPAM:
		return domain.PlaceReportReasonSpam
	case placesv1.PlaceReportReason_PLACE_REPORT_REASON_OFFENSIVE_CONTENT:
		return domain.PlaceReportReasonOffensiveContent
	case placesv1.PlaceReportReason_PLACE_REPORT_REASON_WRONG_INFO:
		return domain.PlaceReportReasonWrongInfo
	case placesv1.PlaceReportReason_PLACE_REPORT_REASON_DUPLICATE:
		return domain.PlaceReportReasonDuplicate
	case placesv1.PlaceReportReason_PLACE_REPORT_REASON_CLOSED_PLACE:
		return domain.PlaceReportReasonClosedPlace
	case placesv1.PlaceReportReason_PLACE_REPORT_REASON_OTHER:
		return domain.PlaceReportReasonOther
	default:
		return ""
	}
}

func reportReasonToProto(reason domain.PlaceReportReason) placesv1.PlaceReportReason {
	switch reason {
	case domain.PlaceReportReasonSpam:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_SPAM
	case domain.PlaceReportReasonOffensiveContent:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_OFFENSIVE_CONTENT
	case domain.PlaceReportReasonWrongInfo:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_WRONG_INFO
	case domain.PlaceReportReasonDuplicate:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_DUPLICATE
	case domain.PlaceReportReasonClosedPlace:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_CLOSED_PLACE
	case domain.PlaceReportReasonOther:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_OTHER
	default:
		return placesv1.PlaceReportReason_PLACE_REPORT_REASON_UNSPECIFIED
	}
}

func placeReportToProto(report domain.PlaceReport) *placesv1.PlaceReport {
	protoReport := &placesv1.PlaceReport{
		Id:               report.ID,
		PlaceId:          report.PlaceID,
		ReportedByUserId: report.ReportedByUserID,
		Reason:           reportReasonToProto(report.Reason),
		Comment:          report.Comment,
		Status:           reportStatusToProto(report.Status),
		CreatedAt:        timestamppb.New(report.CreatedAt),
	}

	if report.ResolvedAt != nil {
		protoReport.ResolvedAt = timestamppb.New(*report.ResolvedAt)
	}

	return protoReport
}

func reportStatusToProto(status domain.PlaceReportStatus) placesv1.PlaceReportStatus {
	switch status {
	case domain.PlaceReportStatusOpen:
		return placesv1.PlaceReportStatus_PLACE_REPORT_STATUS_OPEN
	case domain.PlaceReportStatusResolved:
		return placesv1.PlaceReportStatus_PLACE_REPORT_STATUS_RESOLVED
	case domain.PlaceReportStatusRejected:
		return placesv1.PlaceReportStatus_PLACE_REPORT_STATUS_REJECTED
	default:
		return placesv1.PlaceReportStatus_PLACE_REPORT_STATUS_UNSPECIFIED
	}
}

func categoriesToProto(categories []domain.Category) []*placesv1.Category {
	protoCategories := make([]*placesv1.Category, len(categories))

	for i, category := range categories {
		protoCategories[i] = categoryToProto(category)
	}

	return protoCategories
}

func categoryToProto(category domain.Category) *placesv1.Category {
	return &placesv1.Category{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: timestamppb.New(category.CreatedAt),
	}
}
