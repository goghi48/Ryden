package storage

import (
	"context"
	"sync"

	"github.com/goghi48/ryden/internal/place/domain"
)

type MemoryStorage struct {
	places  map[string]domain.Place
	reports map[string]domain.PlaceReport
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		places:  make(map[string]domain.Place),
		reports: make(map[string]domain.PlaceReport),
	}
}

func (s *MemoryStorage) Create(ctx context.Context, place domain.Place, categoryIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.places[place.ID]; ok {
		return ErrPlaceAlreadyExists
	}

	s.places[place.ID] = place
	return nil
}

func (s *MemoryStorage) GetByID(ctx context.Context, id string) (domain.Place, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	place, ok := s.places[id]
	if !ok {
		return domain.Place{}, ErrPlaceNotFound
	}

	return place, nil
}

func (s *MemoryStorage) List(ctx context.Context, city string, limit int) ([]domain.Place, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var places []domain.Place
	for _, place := range s.places {
		if city != "" && place.City != city {
			continue
		}

		if place.Status != domain.StatusApproved {
			continue
		}

		places = append(places, place)

		if limit > 0 && len(places) >= limit {
			break
		}
	}

	return places, nil
}

func (s *MemoryStorage) CreateReport(ctx context.Context, report domain.PlaceReport, reviewThreshold int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	place, ok := s.places[report.PlaceID]
	if !ok {
		return ErrPlaceNotFound
	}

	if _, ok := s.reports[report.ID]; ok {
		return ErrPlaceReportAlreadyExists
	}

	for _, existingReport := range s.reports {
		if existingReport.PlaceID == report.PlaceID &&
			existingReport.ReportedByUserID == report.ReportedByUserID {
			return ErrPlaceReportAlreadyExists
		}
	}

	s.reports[report.ID] = report

	if reviewThreshold > 0 && countOpenReportsByPlaceID(s.reports, report.PlaceID) >= reviewThreshold {
		place.Status = domain.StatusPendingReview
		place.UpdatedAt = report.CreatedAt
		s.places[report.PlaceID] = place
	}

	return nil
}

func countOpenReportsByPlaceID(reports map[string]domain.PlaceReport, placeID string) int {
	count := 0

	for _, report := range reports {
		if report.PlaceID == placeID && report.Status == domain.PlaceReportStatusOpen {
			count++
		}
	}

	return count
}
