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

func (s *PostgresStorage) Create(ctx context.Context, place domain.Place, categoryIDs []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const insertPlaceQuery = `
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

	_, err = tx.Exec(
		ctx,
		insertPlaceQuery,
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
		return mapPostgresCreateError(err)
	}

	const insertPlaceCategoryQuery = `
		INSERT INTO place_categories (
			place_id,
			category_id
		)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	for _, categoryID := range categoryIDs {
		_, err := tx.Exec(
			ctx,
			insertPlaceCategoryQuery,
			place.ID,
			categoryID,
		)
		if err != nil {
			return mapPostgresCreateError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
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

	categories, err := s.getCategoriesByPlaceID(ctx, place.ID)
	if err != nil {
		return domain.Place{}, err
	}

	place.Categories = categories

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
		AND status = 'APPROVED'
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

	if len(places) == 0 {
		return places, nil
	}

	categoriesByPlaceID, err := s.getCategoriesByPlaceIDs(ctx, extractPlaceIDs(places))
	if err != nil {
		return nil, err
	}

	for i := range places {
		places[i].Categories = categoriesByPlaceID[places[i].ID]
	}

	return places, nil
}

func mapPostgresCreateError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		return ErrPlaceAlreadyExists
	case "23503":
		return ErrCategoryNotFound
	default:
		return err
	}
}

func (s *PostgresStorage) getCategoriesByPlaceID(ctx context.Context, placeID string) ([]domain.Category, error) {
	const query = `
		SELECT
			c.id,
			c.name,
			c.slug,
			c.created_at
		FROM categories c
		JOIN place_categories pc ON pc.category_id = c.id
		WHERE pc.place_id = $1
		ORDER BY c.name
	`

	rows, err := s.pool.Query(ctx, query, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)

	for rows.Next() {
		var category domain.Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.CreatedAt,
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func extractPlaceIDs(places []domain.Place) []string {
	ids := make([]string, len(places))

	for i, place := range places {
		ids[i] = place.ID
	}

	return ids
}

func (s *PostgresStorage) getCategoriesByPlaceIDs(ctx context.Context, placeIDs []string) (map[string][]domain.Category, error) {
	const query = `
		SELECT
			pc.place_id,
			c.id,
			c.name,
			c.slug,
			c.created_at
		FROM place_categories pc
		JOIN categories c ON c.id = pc.category_id
		WHERE pc.place_id = ANY($1::uuid[])
		ORDER BY c.name
	`

	rows, err := s.pool.Query(ctx, query, placeIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoriesByPlaceID := make(map[string][]domain.Category)

	for rows.Next() {
		var placeID string
		var category domain.Category

		if err := rows.Scan(
			&placeID,
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.CreatedAt,
		); err != nil {
			return nil, err
		}

		categoriesByPlaceID[placeID] = append(categoriesByPlaceID[placeID], category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoriesByPlaceID, nil
}

func (s *PostgresStorage) CreateReport(ctx context.Context, report domain.PlaceReport, reviewThreshold int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const insertReportQuery = `
		INSERT INTO place_reports (
			id,
			place_id,
			reported_by_user_id,
			reason,
			comment,
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.Exec(
		ctx,
		insertReportQuery,
		report.ID,
		report.PlaceID,
		report.ReportedByUserID,
		string(report.Reason),
		report.Comment,
		string(report.Status),
		report.CreatedAt,
	)
	if err != nil {
		return mapPostgresCreateReportError(err)
	}

	if reviewThreshold > 0 {
		var openReportsCount int

		const countReportsQuery = `
			SELECT count(*)
			FROM place_reports
			WHERE place_id = $1
				AND status = 'OPEN'
		`

		err = tx.QueryRow(ctx, countReportsQuery, report.PlaceID).Scan(&openReportsCount)
		if err != nil {
			return err
		}

		if openReportsCount >= reviewThreshold {
			const updatePlaceStatusQuery = `
				UPDATE places
				SET status = $1,
					updated_at = $2
				WHERE id = $3
					AND status = $4
			`

			_, err = tx.Exec(
				ctx,
				updatePlaceStatusQuery,
				string(domain.StatusPendingReview),
				report.CreatedAt,
				report.PlaceID,
				string(domain.StatusApproved),
			)
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func mapPostgresCreateReportError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		return ErrPlaceReportAlreadyExists
	case "23503":
		return ErrPlaceNotFound
	default:
		return err
	}
}
