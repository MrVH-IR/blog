package tags

import (
	"blog/internal/auth"
	"blog/internal/httpx"
	"blog/internal/post"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	}

	var req CreateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		httpx.WriteError(w, http.StatusBadRequest, "request must contain a single json object")
		return
	}

	_, err := post.ParseUserID(userID)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	tags, err := h.Service.Create(r.Context(), req.TaggableID, &req)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTagType):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidTagID):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidTagName):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, tags)
}
