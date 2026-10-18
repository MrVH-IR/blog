package tags

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidTagType = errors.New("invalid tag type")
	ErrInvalidTagID   = errors.New("invalid tag id")
	ErrInvalidTagName = errors.New("invalid tag name")
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(ctx context.Context, typeID int64, req *CreateRequest) (*Tags, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.TaggableType = strings.TrimSpace(req.TaggableType)

	switch {
	case req.Name == "":
		return nil, ErrInvalidTagName
	case req.TaggableType == "":
		return nil, ErrInvalidTagType
	case req.TaggableID == 0:
		return nil, ErrInvalidTagID
	}

	return s.Repo.Create(ctx, typeID, req)
}
