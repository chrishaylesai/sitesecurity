package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrishaylesai/sitesecurity/api/internal/auth"
	"github.com/chrishaylesai/sitesecurity/api/internal/middleware"
)

// mockProvider is a test double for auth.Provider.
type mockProvider struct {
	validateFunc func(ctx context.Context, token string) (*auth.Claims, error)
}

func (m *mockProvider) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	return m.validateFunc(ctx, token)
}

func (m *mockProvider) GetUserInfo(ctx context.Context, token string) (*auth.UserInfo, error) {
	return nil, nil
}

func (m *mockProvider) GetLoginURL(state string) string {
	return ""
}

func (m *mockProvider) ExchangeCode(ctx context.Context, code string) (*auth.TokenSet, error) {
	return nil, nil
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	provider := &mockProvider{}
	mw := middleware.Auth(provider)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	provider := &mockProvider{
		validateFunc: func(ctx context.Context, token string) (*auth.Claims, error) {
			return nil, auth.ErrInvalidToken
		},
	}
	mw := middleware.Auth(provider)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	claims := &auth.Claims{
		Subject: "user-123",
		Email:   "test@example.com",
		Name:    "Test User",
		Roles:   []string{"worker"},
	}

	provider := &mockProvider{
		validateFunc: func(ctx context.Context, token string) (*auth.Claims, error) {
			return claims, nil
		},
	}
	mw := middleware.Auth(provider)

	called := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		gotClaims := middleware.GetClaims(r.Context())
		if gotClaims == nil {
			t.Error("expected claims in context")
			return
		}
		if gotClaims.Subject != "user-123" {
			t.Errorf("expected subject 'user-123', got '%s'", gotClaims.Subject)
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("handler was not called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRequireRole_HasRole(t *testing.T) {
	claims := &auth.Claims{
		Subject: "user-123",
		Roles:   []string{"company_admin"},
	}

	mw := middleware.RequireRole("company_admin")
	called := false

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), middleware.ClaimsContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("handler was not called")
	}
}

func TestRequireRole_MissingRole(t *testing.T) {
	claims := &auth.Claims{
		Subject: "user-123",
		Roles:   []string{"worker"},
	}

	mw := middleware.RequireRole("company_admin")

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), middleware.ClaimsContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}

func TestRequireRole_NoClaims(t *testing.T) {
	mw := middleware.RequireRole("worker")

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
