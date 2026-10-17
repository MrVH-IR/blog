package auth

import (
	"blog/internal/httpx"
	"blog/internal/user"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type LoginService interface {
	Login(ctx context.Context, req *LoginRequest) (*user.User, error)
}

type Handler struct {
	Service LoginService
}

func NewHandler(s LoginService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "Cannot decode request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	u, err := h.Service.Login(ctx, &req)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpx.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"id":         strconv.Itoa(u.ID),
		"name":       u.Name,
		"email":      u.Email,
		"created_at": u.CreatedAt.Format(time.RFC3339),
	})
}
