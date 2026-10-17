package auth

import (
	"blog/internal/user"
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Service struct {
	UserRepo *user.Repository
}

func NewService(userRepo *user.Repository) *Service {
	return &Service{UserRepo: userRepo}
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*user.User, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := s.UserRepo.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}
