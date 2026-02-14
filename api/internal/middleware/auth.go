package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chrishaylesai/sitesecurity/api/internal/auth"
)

type contextKey string

const ClaimsContextKey contextKey = "auth_claims"

// Auth returns middleware that validates the Authorization header
// and injects claims into the request context.
func Auth(provider auth.Provider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				http.Error(w, `{"error": "missing or invalid authorization header"}`, http.StatusUnauthorized)
				return
			}

			claims, err := provider.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims extracts auth claims from the request context.
func GetClaims(ctx context.Context) *auth.Claims {
	claims, ok := ctx.Value(ClaimsContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

// RequireRole returns middleware that checks the user has one of the specified roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
				return
			}

			for _, required := range roles {
				for _, userRole := range claims.Roles {
					if userRole == required {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
			return
		})
	}
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return parts[1]
}
