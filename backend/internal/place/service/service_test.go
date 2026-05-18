package service

import (
	"context"
	"errors"
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
		CreatedByUserID: "user-1",
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

	if place.Status != domain.StatusPendingReview {
		t.Fatalf("expected status %s, got %s", domain.StatusPendingReview, place.Status)
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
