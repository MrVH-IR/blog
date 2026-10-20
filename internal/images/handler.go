package images

import (
	"errors"
	"net/http"
	"strconv"

	"blog/internal/httpx"
)

const maxImageSize = 5 << 20 // 5 MiB

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Limit the entire request body to 5 MiB.
	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize)

	postID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || postID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	// Parse multipart/form-data.
	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid multipart request or file too large")
		return
	}

	// Get the uploaded file from the "image" field.
	file, _, err := r.FormFile("image")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	alt := r.FormValue("alt")

	image, err := h.Service.Create(
		r.Context(),
		postID,
		file,
		alt,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidImage):
			httpx.WriteError(w, http.StatusBadRequest, "unsupported image format")
		case errors.Is(err, ErrInvalidAlt):
			httpx.WriteError(w, http.StatusBadRequest, "alt must be at most 255 characters")
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "failed to save image")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, image)
}
