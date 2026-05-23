package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/goghi48/ryden/internal/place/domain"
	"github.com/goghi48/ryden/internal/place/storage"
)

func newTestPlaceService(t *testing.T) *PlaceService {
	t.Helper()

	memoryStorage := storage.NewMemoryStorage()
	return NewPlaceService(memoryStorage)
}

func validCreatePlaceInput() CreatePlaceInput {
	return CreatePlaceInput{
		Title:           "test title",
		Description:     "test description",
		Address:         "test address",
		City:            "Moscow",
		Latitude:        55.7558,
		Longitude:       37.6173,
		CreatedByUserID: "11111111-1111-1111-1111-111111111111",
		CategoryIDs: []string{
			"11111111-1111-1111-1111-000000000002",
		},
	}
}

func validCreatePlaceReportInput(placeID string) CreatePlaceReportInput {
	return CreatePlaceReportInput{
		PlaceID:          placeID,
		ReportedByUserID: "22222222-2222-2222-2222-222222222222",
		Reason:           domain.PlaceReportReasonWrongInfo,
		Comment:          "wrong address",
	}
}

func mustCreatePlace(t *testing.T, placeService *PlaceService, input CreatePlaceInput) domain.Place {
	t.Helper()

	place, err := placeService.CreatePlace(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error while creating place, got %v", err)
	}

	return place
}

func TestPlaceService_CreatePlace_Success(t *testing.T) {
	placeService := newTestPlaceService(t)
	input := validCreatePlaceInput()

	place, err := placeService.CreatePlace(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if place.ID == "" {
		t.Fatal("expected place ID to be generated")
	}

	if place.Status != domain.StatusApproved {
		t.Fatalf("expected status %s, got %s", domain.StatusApproved, place.Status)
	}

	if place.Title != input.Title {
		t.Fatalf("expected title %q, got %q", input.Title, place.Title)
	}

	if place.City != input.City {
		t.Fatalf("expected city %q, got %q", input.City, place.City)
	}

	if place.CreatedByUserID != input.CreatedByUserID {
		t.Fatalf("expected created by user ID %q, got %q", input.CreatedByUserID, place.CreatedByUserID)
	}

	if place.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	if place.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
}

func TestPlaceService_CreatePlace_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		mutateInput func(input *CreatePlaceInput)
		expectedErr error
	}{
		{
			name: "without title",
			mutateInput: func(input *CreatePlaceInput) {
				input.Title = ""
			},
			expectedErr: ErrTitleRequired,
		},
		{
			name: "without city",
			mutateInput: func(input *CreatePlaceInput) {
				input.City = ""
			},
			expectedErr: ErrCityRequired,
		},
		{
			name: "invalid latitude",
			mutateInput: func(input *CreatePlaceInput) {
				input.Latitude = 91
			},
			expectedErr: ErrInvalidLatitude,
		},
		{
			name: "invalid longitude",
			mutateInput: func(input *CreatePlaceInput) {
				input.Longitude = 181
			},
			expectedErr: ErrInvalidLongitude,
		},
		{
			name: "invalid created by user id",
			mutateInput: func(input *CreatePlaceInput) {
				input.CreatedByUserID = "user-1"
			},
			expectedErr: ErrInvalidCreatedByUserID,
		},
		{
			name: "invalid category id",
			mutateInput: func(input *CreatePlaceInput) {
				input.CategoryIDs = []string{"bad-category-id"}
			},
			expectedErr: ErrInvalidCategoryID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			placeService := newTestPlaceService(t)
			input := validCreatePlaceInput()

			tt.mutateInput(&input)

			_, err := placeService.CreatePlace(context.Background(), input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestPlaceService_GetPlace_Success(t *testing.T) {
	placeService := newTestPlaceService(t)
	input := validCreatePlaceInput()

	createdPlace := mustCreatePlace(t, placeService, input)

	foundPlace, err := placeService.GetPlace(context.Background(), createdPlace.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if foundPlace.ID != createdPlace.ID {
		t.Fatalf("expected place ID %q, got %q", createdPlace.ID, foundPlace.ID)
	}

	if foundPlace.Title != createdPlace.Title {
		t.Fatalf("expected title %q, got %q", createdPlace.Title, foundPlace.Title)
	}
}

func TestPlaceService_GetPlace_NotFound(t *testing.T) {
	placeService := newTestPlaceService(t)

	_, err := placeService.GetPlace(context.Background(), "not-existing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, storage.ErrPlaceNotFound) {
		t.Fatalf("expected error %v, got %v", storage.ErrPlaceNotFound, err)
	}
}

func TestPlaceService_ListPlaces_FiltersByCity(t *testing.T) {
	placeService := newTestPlaceService(t)

	moscowInput := validCreatePlaceInput()
	moscowInput.City = "Moscow"

	novosibirskInput := validCreatePlaceInput()
	novosibirskInput.City = "Novosibirsk"
	novosibirskInput.Latitude = 55.0084
	novosibirskInput.Longitude = 82.9357

	mustCreatePlace(t, placeService, moscowInput)
	mustCreatePlace(t, placeService, novosibirskInput)

	places, err := placeService.ListPlaces(context.Background(), "Moscow", 10)
	if err != nil {
		t.Fatalf("expected no error while listing places, got %v", err)
	}

	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}

	if places[0].City != "Moscow" {
		t.Fatalf("expected place city Moscow, got %s", places[0].City)
	}
}

func TestPlaceService_ListPlaces_RespectsLimit(t *testing.T) {
	placeService := newTestPlaceService(t)

	firstInput := validCreatePlaceInput()
	firstInput.Title = "first place"

	secondInput := validCreatePlaceInput()
	secondInput.Title = "second place"
	secondInput.Latitude = 55.7600
	secondInput.Longitude = 37.6200

	mustCreatePlace(t, placeService, firstInput)
	mustCreatePlace(t, placeService, secondInput)

	places, err := placeService.ListPlaces(context.Background(), "", 1)
	if err != nil {
		t.Fatalf("expected no error while listing places, got %v", err)
	}

	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}

func TestPlaceService_CreatePlaceReport_Success(t *testing.T) {
	placeService := newTestPlaceService(t)
	place := mustCreatePlace(t, placeService, validCreatePlaceInput())
	input := validCreatePlaceReportInput(place.ID)

	report, err := placeService.CreatePlaceReport(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.ID == "" {
		t.Fatal("expected report ID to be generated")
	}

	if report.PlaceID != input.PlaceID {
		t.Fatalf("expected place ID %q, got %q", input.PlaceID, report.PlaceID)
	}

	if report.ReportedByUserID != input.ReportedByUserID {
		t.Fatalf("expected reported by user ID %q, got %q", input.ReportedByUserID, report.ReportedByUserID)
	}

	if report.Reason != input.Reason {
		t.Fatalf("expected reason %q, got %q", input.Reason, report.Reason)
	}

	if report.Comment != input.Comment {
		t.Fatalf("expected comment %q, got %q", input.Comment, report.Comment)
	}

	if report.Status != domain.PlaceReportStatusOpen {
		t.Fatalf("expected status %q, got %q", domain.PlaceReportStatusOpen, report.Status)
	}

	if report.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	if report.ResolvedAt != nil {
		t.Fatal("expected ResolvedAt to be nil")
	}
}

func TestPlaceService_CreatePlaceReport_ValidationErrors(t *testing.T) {
	placeService := newTestPlaceService(t)
	place := mustCreatePlace(t, placeService, validCreatePlaceInput())

	tests := []struct {
		name        string
		mutateInput func(input *CreatePlaceReportInput)
		expectedErr error
	}{
		{
			name: "invalid place id",
			mutateInput: func(input *CreatePlaceReportInput) {
				input.PlaceID = "bad-place-id"
			},
			expectedErr: ErrInvalidPlaceID,
		},
		{
			name: "invalid reported by user id",
			mutateInput: func(input *CreatePlaceReportInput) {
				input.ReportedByUserID = "bad-user-id"
			},
			expectedErr: ErrInvalidReportedByUserID,
		},
		{
			name: "invalid reason",
			mutateInput: func(input *CreatePlaceReportInput) {
				input.Reason = ""
			},
			expectedErr: ErrInvalidReportReason,
		},
		{
			name: "comment too long",
			mutateInput: func(input *CreatePlaceReportInput) {
				input.Comment = strings.Repeat("a", 501)
			},
			expectedErr: ErrReportCommentTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validCreatePlaceReportInput(place.ID)
			tt.mutateInput(&input)

			_, err := placeService.CreatePlaceReport(context.Background(), input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestPlaceService_CreatePlaceReport_PlaceNotFound(t *testing.T) {
	placeService := newTestPlaceService(t)
	input := validCreatePlaceReportInput("99999999-9999-9999-9999-999999999999")

	_, err := placeService.CreatePlaceReport(context.Background(), input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, storage.ErrPlaceNotFound) {
		t.Fatalf("expected error %v, got %v", storage.ErrPlaceNotFound, err)
	}
}

func TestPlaceService_CreatePlaceReport_DuplicateUserPerPlace(t *testing.T) {
	placeService := newTestPlaceService(t)
	place := mustCreatePlace(t, placeService, validCreatePlaceInput())
	input := validCreatePlaceReportInput(place.ID)

	if _, err := placeService.CreatePlaceReport(context.Background(), input); err != nil {
		t.Fatalf("expected no error while creating first report, got %v", err)
	}

	_, err := placeService.CreatePlaceReport(context.Background(), input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, storage.ErrPlaceReportAlreadyExists) {
		t.Fatalf("expected error %v, got %v", storage.ErrPlaceReportAlreadyExists, err)
	}
}

func TestPlaceService_CreatePlaceReport_MarksPlacePendingReviewWhenThresholdReached(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	placeService := NewPlaceServiceWithReportsReviewThreshold(memoryStorage, 2)

	place := mustCreatePlace(t, placeService, validCreatePlaceInput())

	firstReportInput := validCreatePlaceReportInput(place.ID)
	if _, err := placeService.CreatePlaceReport(context.Background(), firstReportInput); err != nil {
		t.Fatalf("expected no error while creating first report, got %v", err)
	}

	foundPlace, err := placeService.GetPlace(context.Background(), place.ID)
	if err != nil {
		t.Fatalf("expected no error while getting place, got %v", err)
	}
	if foundPlace.Status != domain.StatusApproved {
		t.Fatalf("expected status %q before threshold, got %q", domain.StatusApproved, foundPlace.Status)
	}

	secondReportInput := validCreatePlaceReportInput(place.ID)
	secondReportInput.ReportedByUserID = "33333333-3333-3333-3333-333333333333"

	if _, err := placeService.CreatePlaceReport(context.Background(), secondReportInput); err != nil {
		t.Fatalf("expected no error while creating second report, got %v", err)
	}

	foundPlace, err = placeService.GetPlace(context.Background(), place.ID)
	if err != nil {
		t.Fatalf("expected no error while getting place, got %v", err)
	}
	if foundPlace.Status != domain.StatusPendingReview {
		t.Fatalf("expected status %q after threshold, got %q", domain.StatusPendingReview, foundPlace.Status)
	}
}

func TestPlaceService_CreatePlaceReport_DoesNotMarkPlacePendingReviewWhenAutoReviewDisabled(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	placeService := NewPlaceServiceWithReportsAutoReview(memoryStorage, false, 1)

	place := mustCreatePlace(t, placeService, validCreatePlaceInput())

	reportInput := validCreatePlaceReportInput(place.ID)
	if _, err := placeService.CreatePlaceReport(context.Background(), reportInput); err != nil {
		t.Fatalf("expected no error while creating report, got %v", err)
	}

	foundPlace, err := placeService.GetPlace(context.Background(), place.ID)
	if err != nil {
		t.Fatalf("expected no error while getting place, got %v", err)
	}

	if foundPlace.Status != domain.StatusApproved {
		t.Fatalf("expected status %q when auto review disabled, got %q", domain.StatusApproved, foundPlace.Status)
	}
}
