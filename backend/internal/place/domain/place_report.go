package domain

import "time"

type PlaceReportReason string

const (
	PlaceReportReasonSpam             PlaceReportReason = "SPAM"
	PlaceReportReasonOffensiveContent PlaceReportReason = "OFFENSIVE_CONTENT"
	PlaceReportReasonWrongInfo        PlaceReportReason = "WRONG_INFO"
	PlaceReportReasonDuplicate        PlaceReportReason = "DUPLICATE"
	PlaceReportReasonClosedPlace      PlaceReportReason = "CLOSED_PLACE"
	PlaceReportReasonOther            PlaceReportReason = "OTHER"
)

type PlaceReportStatus string

const (
	PlaceReportStatusOpen     PlaceReportStatus = "OPEN"
	PlaceReportStatusResolved PlaceReportStatus = "RESOLVED"
	PlaceReportStatusRejected PlaceReportStatus = "REJECTED"
)

type PlaceReport struct {
	ID               string
	PlaceID          string
	ReportedByUserID string
	Reason           PlaceReportReason
	Comment          string
	Status           PlaceReportStatus
	CreatedAt        time.Time
	ResolvedAt       *time.Time
}
