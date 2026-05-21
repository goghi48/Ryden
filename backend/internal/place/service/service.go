package service

import (
	"context"
	"strings"
	"time"

	"github.com/goghi48/ryden/internal/place/domain"
	"github.com/google/uuid"
)

type PlaceStorage interface {
	Create(ctx context.Context, place domain.Place, categoryIDs []string) error
	GetByID(ctx context.Context, id string) (domain.Place, error)
	List(ctx context.Context, city string, limit int) ([]domain.Place, error)
}

type PlaceService struct {
	storage PlaceStorage
}

func NewPlaceService(storage PlaceStorage) *PlaceService {
	return &PlaceService{storage: storage}
}

func (s *PlaceService) CreatePlace(ctx context.Context, input CreatePlaceInput) (domain.Place, error) {
	if err := validateCreatePlaceInput(input); err != nil {
		return domain.Place{}, err
	}

	id := uuid.NewString()
	now := time.Now()

	place := domain.Place{
		ID:              id,
		Title:           input.Title,
		Description:     input.Description,
		Address:         input.Address,
		City:            input.City,
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		CreatedByUserID: input.CreatedByUserID,
		Status:          domain.StatusApproved,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := s.storage.Create(ctx, place, input.CategoryIDs)
	if err != nil {
		return domain.Place{}, err
	}

	return place, nil
}

func (s *PlaceService) GetPlace(ctx context.Context, id string) (domain.Place, error) {
	place, err := s.storage.GetByID(ctx, id)

	if err != nil {
		return domain.Place{}, err
	}
	return place, nil
}

func (s *PlaceService) ListPlaces(ctx context.Context, city string, limit int) ([]domain.Place, error) {
	places, err := s.storage.List(ctx, city, limit)

	if err != nil {
		return nil, err
	}
	return places, nil
}

func validateCreatePlaceInput(input CreatePlaceInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return ErrTitleRequired
	}
	if strings.TrimSpace(input.City) == "" {
		return ErrCityRequired
	}
	if input.Latitude < -90 || input.Latitude > 90 {
		return ErrInvalidLatitude
	}
	if input.Longitude < -180 || input.Longitude > 180 {
		return ErrInvalidLongitude
	}
	if _, err := uuid.Parse(input.CreatedByUserID); err != nil {
		return ErrInvalidCreatedByUserID
	}

	for _, categoryID := range input.CategoryIDs {
		if _, err := uuid.Parse(categoryID); err != nil {
			return ErrInvalidCategoryID
		}
	}

	return nil
}

type CreatePlaceInput struct {
	Title           string
	Description     string
	Address         string
	City            string
	Latitude        float64
	Longitude       float64
	CreatedByUserID string
	CategoryIDs     []string
}
