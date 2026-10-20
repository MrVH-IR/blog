package images

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

func (r *Repository) Create(ctx context.Context, req *CreateRequest) (*Image, error) {
	query := `INSERT INTO images (imageable_type,imageable_id,path,alt) VALUES ($1, $2, $3, $4) RETURNING id,imageable_type,imageable_id,path,alt,created_at`

	image := &Image{}

	err := r.DB.QueryRow(
		ctx,
		query,
		req.ImageableType,
		req.ImageableID,
		req.Path,
		req.Alt,
	).Scan(
		&image.ID,
		&image.ImageableType,
		&image.ImageableID,
		&image.Path,
		&image.Alt,
		&image.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return image, nil
}
