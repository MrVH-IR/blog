package tags

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, typeID int64, req *CreateRequest) (*Tags, error) {
	query := `INSERT INTO tags (name, taggable_type, taggable_id) VALUES ($1, $2, $3) RETURNING id, name, taggable_type, taggable_id, created_at`

	tags := &Tags{}

	err := r.DB.QueryRow(ctx, query, req.Name, req.TaggableType, typeID).Scan(&tags.ID, &tags.Name, &tags.TaggableType, &tags.TaggableID)

	if err != nil {
		return nil, err
	}

	return tags, nil
}
