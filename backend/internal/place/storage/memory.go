package storage

import (
	"sync"

	"github.com/goghi48/ryden/internal/place/domain"
)

type MemoryStorage struct {
	places map[string]domain.Place
	mu     sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		places: make(map[string]domain.Place),
	}
}

func (s *MemoryStorage) Create(place domain.Place) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.places[place.ID]; ok {
		return ErrPlaceAlreadyExists
	}

	s.places[place.ID] = place
	return nil
}

func (s *MemoryStorage) GetByID(id string) (domain.Place, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	place, ok := s.places[id]
	if !ok {
		return domain.Place{}, ErrPlaceNotFound
	}

	return place, nil
}

func (s *MemoryStorage) List(city string, limit int) ([]domain.Place, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var places []domain.Place
	for _, place := range s.places {
		if city != "" && place.City != city {
			continue
		}

		places = append(places, place)

		if limit > 0 && len(places) >= limit {
			break
		}
	}

	return places, nil
}
