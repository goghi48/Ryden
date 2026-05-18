package storage

import (
	"context"
	"errors"

	"github.com/goghi48/ryden/internal/place/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		pool: pool,
	}
}

func (s *PostgresStorage) Create(ctx context.Context, place domain.Place) error {
	const query = `
		INSERT INTO places (
			id,
			title,
			description,
			address,
			city,
			latitude,
			longitude,
			created_by_user_id,
			status,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11
		)
	`

	_, err := s.pool.Exec(
		ctx,
		query,
		place.ID,
		place.Title,
		place.Description,
		place.Address,
		place.City,
		place.Latitude,
		place.Longitude,
		place.CreatedByUserID,
		string(place.Status),
		place.CreatedAt,
		place.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrPlaceAlreadyExists
		}

		return err
	}

	return nil
}

func (s *PostgresStorage) GetByID(ctx context.Context, id string) (domain.Place, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			address,
			city,
			latitude,
			longitude,
			created_by_user_id,
			status,
			created_at,
			updated_at
		FROM places
		WHERE id = $1
	`

	var place domain.Place

	err := s.pool.QueryRow(ctx, query, id).Scan(
		&place.ID,
		&place.Title,
		&place.Description,
		&place.Address,
		&place.City,
		&place.Latitude,
		&place.Longitude,
		&place.CreatedByUserID,
		&place.Status,
		&place.CreatedAt,
		&place.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Place{}, ErrPlaceNotFound
		}

		return domain.Place{}, err
	}
	return place, nil
}

func (s *PostgresStorage) List(ctx context.Context, city string, limit int) ([]domain.Place, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			address,
			city,
			latitude,
			longitude,
			created_by_user_id,
			status,
			created_at,
			updated_at
		FROM places
		WHERE ($1 = '' OR city = $1)
		ORDER BY created_at DESC
		LIMIT CASE WHEN $2::int > 0 THEN $2::int ELSE NULL END
	`

	rows, err := s.pool.Query(ctx, query, city, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	places := make([]domain.Place, 0)

	for rows.Next() {
		var place domain.Place

		if err := rows.Scan(
			&place.ID,
			&place.Title,
			&place.Description,
			&place.Address,
			&place.City,
			&place.Latitude,
			&place.Longitude,
			&place.CreatedByUserID,
			&place.Status,
			&place.CreatedAt,
			&place.UpdatedAt,
		); err != nil {
			return nil, err
		}

		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}
