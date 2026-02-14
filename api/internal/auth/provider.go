package auth

import (
	"context"
	"errors"
)

// ErrInvalidToken is returned when a token is invalid or expired.
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims represents the verified claims from an access token.
type Claims struct {
	Subject string
	Email   string
	Name    string
	Roles   []string
}

// UserInfo represents user profile information from the identity provider.
type UserInfo struct {
	Subject   string
	Email     string
	FirstName string
	LastName  string
	Roles     []string
}

// TokenSet represents a set of tokens returned by the identity provider.
type TokenSet struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int
}

// Provider defines the interface for authentication providers.
// Implementations can be swapped (Keycloak, Auth0, Google, etc.)
// by changing the AUTH_PROVIDER environment variable.
type Provider interface {
	// ValidateToken verifies an access token and returns its claims.
	ValidateToken(ctx context.Context, token string) (*Claims, error)

	// GetUserInfo retrieves user profile information using an access token.
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)

	// GetLoginURL returns the URL to redirect users to for login.
	GetLoginURL(state string) string

	// ExchangeCode exchanges an authorization code for a token set.
	ExchangeCode(ctx context.Context, code string) (*TokenSet, error)
}
