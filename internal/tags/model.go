package tags

import "time"

type Tags struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	TaggableType string    `json:"taggable_type"`
	TaggableID   int64     `json:"taggable_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateRequest struct {
	Name         string `json:"name"`
	TaggableType string `json:"taggable_type"`
	TaggableID   int64  `json:"taggable_id"`
}
