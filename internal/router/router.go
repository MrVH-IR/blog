package router

import (
	"blog/internal/auth"
	"blog/internal/httpx"
	"blog/internal/post"
	"blog/internal/tags"
	"blog/internal/user"
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
			})
			return
		}

		httpx.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	mux.Handle("GET /api/v1/me", auth.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(auth.UserIDContextKey).(string)
			if !ok {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
				return
			}

			httpx.WriteJSON(w, http.StatusOK, map[string]string{
				"user_id": userID,
			})
		}),
	))

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	mux.HandleFunc("POST /api/v1/register", userHandler.Register)

	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)

	postRepo := post.NewRepository(pool)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)
	mux.Handle("POST /api/v1/posts", auth.AuthMiddleware(http.HandlerFunc(postHandler.Create)))
	mux.Handle("GET /api/v1/posts", auth.AuthMiddleware(http.HandlerFunc(postHandler.List)))

	tagsRepo := tags.NewRepository(pool)
	tagsService := tags.NewService(tagsRepo)
	tagHandler := tags.NewHandler(tagsService)
	mux.Handle("POST /api/v1/tags", auth.AuthMiddleware(http.HandlerFunc(tagHandler.Create)))

	return mux
}
