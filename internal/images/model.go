package images

import "time"

type Image struct {
	ID            int64     `json:"id"`
	ImageableType string    `json:"imageable_type"`
	ImageableID   int64     `json:"imageable_id"`
	Path          string    `json:"path"`
	Alt           string    `json:"alt"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateRequest struct {
	ImageableType string
	ImageableID   int64
	Path          string
	Alt           string
}
