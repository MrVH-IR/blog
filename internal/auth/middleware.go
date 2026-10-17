package auth

import (
	"blog/internal/httpx"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.Fields(header)
		if len(parts) < 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.WriteError(w, http.StatusUnauthorized, "Token Missing Or Invalid")
			return
		}
		claims, err := ValidateToken(parts[1])
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
