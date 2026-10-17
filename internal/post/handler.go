package post

import (
	"blog/internal/auth"
	"blog/internal/httpx"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDContextKey).(string)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		httpx.WriteError(w, http.StatusBadRequest, "request must contain a single JSON object")
		return
	}

	id, err := parseUserID(userID)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	post, err := h.Service.Create(r.Context(), id, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTitle),
			errors.Is(err, ErrInvalidSlug),
			errors.Is(err, ErrInvalidStatus):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case strings.Contains(err.Error(), "slug must not contain whitespace"),
			strings.Contains(err.Error(), "title and slug must be at most 255 bytes"):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case strings.Contains(err.Error(), "duplicate key"):
			httpx.WriteError(w, http.StatusConflict, "slug already exists")
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "failed to create post")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, post)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Service.List(r.Context())

	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list posts")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, posts)
}

func parseUserID(value string) (int, error) {
	var id int
	_, err := fmt.Sscan(value, &id)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user ID")
	}
	return id, nil
}
