package images

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidImage = errors.New("unsupported image format")
	ErrInvalidAlt   = errors.New("alt must be at most 255 characters")
)

const imageDir = "images"

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	postID int64,
	file io.Reader,
	alt string,
) (*Image, error) {
	alt = strings.TrimSpace(alt)

	if len(alt) > 255 {
		return nil, ErrInvalidAlt
	}

	header := make([]byte, 512)
	n, err := io.ReadFull(file, header)

	if err != nil && !errors.Is(err, io.EOF) &&
		!errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}

	header = header[:n]
	contentType := http.DetectContentType(header)

	extensions := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}

	ext, ok := extensions[contentType]
	if !ok {
		return nil, ErrInvalidImage
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	filename := hex.EncodeToString(randomBytes) + ext

	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(imageDir, filename)

	dst, err := os.OpenFile(
		fullPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		return nil, err
	}

	_, copyErr := io.Copy(dst, io.MultiReader(
		bytes.NewReader(header),
		file,
	))

	closeErr := dst.Close()

	if copyErr != nil {
		_ = os.Remove(fullPath)
		return nil, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(fullPath)
		return nil, closeErr
	}

	req := &CreateRequest{
		ImageableType: "post",
		ImageableID:   postID,
		Path:          filepath.ToSlash(fullPath),
		Alt:           alt,
	}

	image, err := s.Repo.Create(ctx, req)
	if err != nil {
		_ = os.Remove(fullPath)
		return nil, err
	}

	return image, nil
}

func contentTypeFromExtension(ext string) string {
	return mime.TypeByExtension(ext)
}
