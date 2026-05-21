package domain

import "time"

type PlaceStatus string

const (
	StatusPendingReview PlaceStatus = "PENDING_REVIEW"
	StatusApproved      PlaceStatus = "APPROVED"
	StatusRejected      PlaceStatus = "REJECTED"
	StatusArchived      PlaceStatus = "ARCHIVED"
)

type Place struct {
	ID              string
	Title           string
	Description     string
	Address         string
	City            string
	Latitude        float64
	Longitude       float64
	CreatedByUserID string
	Status          PlaceStatus
	Categories      []Category
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
