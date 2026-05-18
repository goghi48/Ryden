package storage

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/goghi48/ryden/internal/place/domain"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func newTestPostgresStorage(t *testing.T) *PostgresStorage {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate postgres container: %v", err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping postgres: %v", err)
	}

	runMigrations(t, connStr)

	return NewPostgresStorage(pool)
}

func runMigrations(t *testing.T, databaseURL string) {
	t.Helper()

	migrationsPath, err := filepath.Abs("../../../migrations/place")
	if err != nil {
		t.Fatalf("failed to get migrations path: %v", err)
	}

	sourceURL := fileURL(migrationsPath)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		t.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("failed to run migrations: %v", err)
	}

	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		t.Fatalf("failed to close migration source: %v", sourceErr)
	}
	if databaseErr != nil {
		t.Fatalf("failed to close migration database: %v", databaseErr)
	}
}

func fileURL(path string) string {
	return "file://" + filepath.ToSlash(path)
}

func validPlace() domain.Place {
	now := time.Now()

	place := domain.Place{
		ID:              "11111111-1111-1111-1111-111111111111",
		Title:           "test title",
		Description:     "test description",
		Address:         "test address",
		City:            "Moscow",
		Latitude:        55.7558,
		Longitude:       37.6173,
		CreatedByUserID: "22222222-2222-2222-2222-222222222222",
		Status:          domain.StatusPendingReview,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return place
}

func TestPostgresStorage_Create_Success(t *testing.T) {
	ctx := context.Background()

	storage := newTestPostgresStorage(t)
	place := validPlace()

	err := storage.Create(ctx, place)
	if err != nil {
		t.Fatalf("failed to create place: %v", err)
	}

	foundPlace, err := storage.GetByID(ctx, place.ID)
	if err != nil {
		t.Fatalf("failed to get created place: %v", err)
	}

	if foundPlace.ID != place.ID {
		t.Fatalf("expected place ID %q, got %q", place.ID, foundPlace.ID)
	}

	if foundPlace.Title != place.Title {
		t.Fatalf("expected title %q, got %q", place.Title, foundPlace.Title)
	}

	if foundPlace.City != place.City {
		t.Fatalf("expected city %q, got %q", place.City, foundPlace.City)
	}

	if foundPlace.Status != place.Status {
		t.Fatalf("expected status %q, got %q", place.Status, foundPlace.Status)
	}

	if foundPlace.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	if foundPlace.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
}

func TestPostgresStorage_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	storage := newTestPostgresStorage(t)

	_, err := storage.GetByID(ctx, "99999999-9999-9999-9999-999999999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrPlaceNotFound) {
		t.Fatalf("expected error %v, got %v", ErrPlaceNotFound, err)
	}
}

func TestPostgresStorage_Create_DuplicateID(t *testing.T) {
	ctx := context.Background()
	storage := newTestPostgresStorage(t)
	place := validPlace()

	if err := storage.Create(ctx, place); err != nil {
		t.Fatalf("expected no error while creating place, got %v", err)
	}

	err := storage.Create(ctx, place)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrPlaceAlreadyExists) {
		t.Fatalf("expected error %v, got %v", ErrPlaceAlreadyExists, err)
	}
}

func TestPostgresStorage_List_FiltersByCity(t *testing.T) {
	ctx := context.Background()
	storage := newTestPostgresStorage(t)

	firstPlace := validPlace()
	firstPlace.City = "Moscow"

	secondPlace := validPlace()
	secondPlace.ID = "33333333-3333-3333-3333-333333333333"
	secondPlace.Title = "second place"
	secondPlace.City = "Novosibirsk"
	secondPlace.Latitude = 55.0084
	secondPlace.Longitude = 82.9357

	if err := storage.Create(ctx, firstPlace); err != nil {
		t.Fatalf("expected no error while creating first place, got %v", err)
	}

	if err := storage.Create(ctx, secondPlace); err != nil {
		t.Fatalf("expected no error while creating second place, got %v", err)
	}

	places, err := storage.List(ctx, "Moscow", 10)
	if err != nil {
		t.Fatalf("expected no error while listing places, got %v", err)
	}

	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}

	if places[0].City != "Moscow" {
		t.Fatalf("expected city Moscow, got %s", places[0].City)
	}
}

func TestPostgresStorage_List_RespectsLimit(t *testing.T) {
	ctx := context.Background()
	storage := newTestPostgresStorage(t)

	firstPlace := validPlace()
	firstPlace.ID = "11111111-1111-1111-1111-111111111111"
	firstPlace.Title = "first place"

	secondPlace := validPlace()
	secondPlace.ID = "33333333-3333-3333-3333-333333333333"
	secondPlace.Title = "second place"
	secondPlace.City = "Novosibirsk"
	secondPlace.Latitude = 55.0084
	secondPlace.Longitude = 82.9357

	if err := storage.Create(ctx, firstPlace); err != nil {
		t.Fatalf("expected no error while creating first place, got %v", err)
	}

	if err := storage.Create(ctx, secondPlace); err != nil {
		t.Fatalf("expected no error while creating second place, got %v", err)
	}

	places, err := storage.List(ctx, "", 1)
	if err != nil {
		t.Fatalf("expected no error while listing places, got %v", err)
	}

	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}
