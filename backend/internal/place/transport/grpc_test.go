package transport

import (
	"context"
	"testing"

	placesv1 "github.com/goghi48/ryden/gen/go/ryden/places/v1"
	"github.com/goghi48/ryden/internal/place/service"
	"github.com/goghi48/ryden/internal/place/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()

	placeStorage := storage.NewMemoryStorage()
	placeService := service.NewPlaceService(placeStorage)
	return NewHandler(placeService)
}

func mustCreatePlaceViaHandler(t *testing.T, handler *Handler) *placesv1.Place {
	t.Helper()

	resp, err := handler.CreatePlace(context.Background(), &placesv1.CreatePlaceRequest{
		Title:           "test title",
		Description:     "test description",
		Address:         "test address",
		City:            "Moscow",
		Latitude:        55.7558,
		Longitude:       37.6173,
		CreatedByUserId: "11111111-1111-1111-1111-111111111111",
	})
	if err != nil {
		t.Fatalf("expected no error while creating place, got %v", err)
	}

	return resp.GetPlace()
}

func TestHandler_CreatePlaceReport_Success(t *testing.T) {
	handler := newTestHandler(t)
	place := mustCreatePlaceViaHandler(t, handler)

	resp, err := handler.CreatePlaceReport(context.Background(), &placesv1.CreatePlaceReportRequest{
		PlaceId:          place.GetId(),
		ReportedByUserId: "22222222-2222-2222-2222-222222222222",
		Reason:           placesv1.PlaceReportReason_PLACE_REPORT_REASON_WRONG_INFO,
		Comment:          "wrong address",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	report := resp.GetPlaceReport()
	if report.GetId() == "" {
		t.Fatal("expected report ID to be generated")
	}
	if report.GetPlaceId() != place.GetId() {
		t.Fatalf("expected place ID %q, got %q", place.GetId(), report.GetPlaceId())
	}
	if report.GetStatus() != placesv1.PlaceReportStatus_PLACE_REPORT_STATUS_OPEN {
		t.Fatalf("expected status OPEN, got %s", report.GetStatus())
	}
	if report.GetReason() != placesv1.PlaceReportReason_PLACE_REPORT_REASON_WRONG_INFO {
		t.Fatalf("expected reason WRONG_INFO, got %s", report.GetReason())
	}
	if report.GetResolvedAt() != nil {
		t.Fatal("expected resolved_at to be nil")
	}
}

func TestHandler_CreatePlaceReport_DuplicateUserPerPlace(t *testing.T) {
	handler := newTestHandler(t)
	place := mustCreatePlaceViaHandler(t, handler)

	req := &placesv1.CreatePlaceReportRequest{
		PlaceId:          place.GetId(),
		ReportedByUserId: "22222222-2222-2222-2222-222222222222",
		Reason:           placesv1.PlaceReportReason_PLACE_REPORT_REASON_WRONG_INFO,
		Comment:          "wrong address",
	}

	if _, err := handler.CreatePlaceReport(context.Background(), req); err != nil {
		t.Fatalf("expected no error while creating first report, got %v", err)
	}

	_, err := handler.CreatePlaceReport(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("expected code %s, got %s", codes.AlreadyExists, status.Code(err))
	}
}

func TestHandler_ApprovePlace_Success(t *testing.T) {
	handler := newTestHandler(t)
	place := mustCreatePlaceViaHandler(t, handler)

	resp, err := handler.ApprovePlace(context.Background(), &placesv1.ApprovePlaceRequest{
		Id: place.GetId(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.GetPlace().GetStatus() != placesv1.PlaceStatus_PLACE_STATUS_APPROVED {
		t.Fatalf("expected status APPROVED, got %s", resp.GetPlace().GetStatus())
	}
}

func TestHandler_ArchivePlace_Success(t *testing.T) {
	handler := newTestHandler(t)
	place := mustCreatePlaceViaHandler(t, handler)

	resp, err := handler.ArchivePlace(context.Background(), &placesv1.ArchivePlaceRequest{
		Id: place.GetId(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.GetPlace().GetStatus() != placesv1.PlaceStatus_PLACE_STATUS_ARCHIVED {
		t.Fatalf("expected status ARCHIVED, got %s", resp.GetPlace().GetStatus())
	}
}

func TestHandler_ApprovePlace_InvalidID(t *testing.T) {
	handler := newTestHandler(t)

	_, err := handler.ApprovePlace(context.Background(), &placesv1.ApprovePlaceRequest{
		Id: "bad-place-id",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected code %s, got %s", codes.InvalidArgument, status.Code(err))
	}
}
