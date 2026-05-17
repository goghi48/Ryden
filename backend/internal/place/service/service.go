package service

import (
	"strings"
	"time"

	"github.com/goghi48/ryden/internal/place/domain"
	"github.com/google/uuid"
)

type PlaceStorage interface {
	Create(place domain.Place) error
	GetByID(id string) (domain.Place, error)
	List(city string, limit int) ([]domain.Place, error)
}

type PlaceService struct {
	storage PlaceStorage
}

func NewPlaceService(storage PlaceStorage) *PlaceService {
	return &PlaceService{storage: storage}
}

func (s *PlaceService) CreatePlace(input CreatePlaceInput) (domain.Place, error) {
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
		Status:          domain.StatusPendingReview,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := s.storage.Create(place)
	if err != nil {
		return domain.Place{}, err
	}

	return place, nil
}

func (s *PlaceService) GetPlace(id string) (domain.Place, error) {
	place, err := s.storage.GetByID(id)

	if err != nil {
		return domain.Place{}, err
	}
	return place, nil
}

func (s *PlaceService) ListPlaces(city string, limit int) ([]domain.Place, error) {
	places, err := s.storage.List(city, limit)

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
}
