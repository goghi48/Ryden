package domain

import "time"

type Category struct {
	ID        string
	Name      string
	Slug      string
	CreatedAt time.Time
}
