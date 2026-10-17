package post

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

func (r *Repository) Create(ctx context.Context, userID int, req *CreateRequest) (*Post, error) {
	query := `INSERT INTO posts (user_id, title, slug, content, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, user_id, title, slug, content, status, created_at`

	post := &Post{}

	err := r.DB.QueryRow(
		ctx,
		query,
		userID,
		req.Title,
		req.Slug,
		req.Content,
		req.Status,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Status,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *Repository) List(ctx context.Context) ([]Post, error) {
	query := `
	SELECT id, user_id, title, slug, content, status, created_at FROM posts
	ORDER BY created_at DESC, id DESC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]Post, 0)

	for rows.Next() {
		var p Post

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
