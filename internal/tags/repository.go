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

	err := r.DB.QueryRow(ctx, query, req.Name, req.TaggableType, typeID).Scan(&tags.ID, &tags.Name, &tags.TaggableType, &tags.TaggableID, &tags.CreatedAt)

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Repository) List(ctx context.Context) ([]*Tags, error) {
	rows, err := r.DB.Query(ctx, `
        SELECT id, name, taggable_type, taggable_id, created_at
        FROM tags
        ORDER BY id DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*Tags, 0)

	for rows.Next() {
		var t Tags

		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.TaggableType,
			&t.TaggableID,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		tags = append(tags, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
