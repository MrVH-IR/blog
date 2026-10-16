package user

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput = errors.New("invalid registration input")
	ErrEmailExists  = errors.New("email already exists")
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Name == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	if utf8.RuneCountInString(req.Name) > 120 || len(req.Email) > 255 {
		return nil, ErrInvalidInput
	}

	if len(req.Password) < 8 || len(req.Password) > 32 {
		return nil, ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	u, err := s.Repo.Create(ctx, req.Name, req.Email, string(passwordHash))

	if err != nil {
		var pgErr interface{ SQLState() string }
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return nil, ErrInvalidInput
		}

		return nil, err
	}

	return u, nil
}
