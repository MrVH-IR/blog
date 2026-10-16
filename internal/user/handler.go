package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"blog/internal/httpx"

	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req RegisterRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON request",
		})
		return
	}


	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	u, err := h.Service.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid registration input",
			})
		case errors.Is(err, ErrEmailExists):
			httpx.WriteJSON(w, http.StatusConflict, map[string]string{
				"error": "email already exists",
			})
		default:
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Printf("registration database error: SQLSTATE=%s", pgErr.Code)
			} else {
				log.Printf("registration error: %v", err)
			}

			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, u)
}
