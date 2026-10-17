package post

import (
	"context"
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidTitle  = errors.New("title is required")
	ErrInvalidSlug   = errors.New("slug is required")
	ErrInvalidStatus = errors.New("status must be draft or published")
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(ctx context.Context, userID int, req *CreateRequest) (*Post, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Slug = strings.TrimSpace(req.Slug)
	req.Status = strings.TrimSpace(req.Status)

	switch {
	case req.Title == "":
		return nil, ErrInvalidTitle
	case req.Slug == "":
		return nil, ErrInvalidSlug
	case req.Status == "":
		return nil, ErrInvalidStatus
	case len(req.Title) > 255:
		return nil, ErrInvalidTitle
	case len(req.Slug) > 255:
		return nil, ErrInvalidSlug
	}

	for _, r := range req.Slug {
		if unicode.IsSpace(r) {
			return nil, errors.New("slug must not contain spaces")
		}
	}

	return s.Repo.Create(ctx, userID, req)
}

func (s *Service) List(ctx context.Context) ([]Post, error) {
	return s.Repo.List(ctx)
}
