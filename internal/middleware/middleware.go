package middleware

import (
	"context"
	"net/http"
	"strings"

	"practice-go-crud/internal/domain"
	"practice-go-crud/internal/service"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := m.extractTokenFromHeader(r)
		if token == "" {
			m.writeError(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		user, err := m.authService.ValidateToken(token)
		if err != nil {
			m.writeError(w, http.StatusUnauthorized, "Invalid authorization token")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func (m *AuthMiddleware) RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*domain.User)
			if !ok {
				m.writeError(w, http.StatusInternalServerError, "User context not found")
				return
			}

			if user.Role != role {
				m.writeError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next(w, r)
		})
	}
}

func (m *AuthMiddleware) extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (m *AuthMiddleware) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
