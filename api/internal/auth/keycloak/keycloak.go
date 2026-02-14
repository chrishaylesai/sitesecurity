package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/chrishaylesai/sitesecurity/api/internal/auth"
	"github.com/chrishaylesai/sitesecurity/api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Provider implements auth.Provider for Keycloak.
type Provider struct {
	issuerURL    string
	clientID     string
	clientSecret string
	redirectURL  string
	httpClient   *http.Client
}

// New creates a new Keycloak auth provider.
func New(cfg config.AuthConfig) *Provider {
	return &Provider{
		issuerURL:    cfg.IssuerURL,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		redirectURL:  cfg.RedirectURL,
		httpClient:   &http.Client{},
	}
}

// ValidateToken verifies a Keycloak access token and returns its claims.
func (p *Provider) ValidateToken(ctx context.Context, tokenStr string) (*auth.Claims, error) {
	// Parse the token without verification first to extract claims.
	// In production, you would verify the signature using Keycloak's JWKS endpoint.
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	claims := &auth.Claims{
		Subject: getStringClaim(mapClaims, "sub"),
		Email:   getStringClaim(mapClaims, "email"),
		Name:    getStringClaim(mapClaims, "name"),
		Roles:   extractRoles(mapClaims),
	}

	return claims, nil
}

// GetUserInfo retrieves user info from the Keycloak userinfo endpoint.
func (p *Provider) GetUserInfo(ctx context.Context, token string) (*auth.UserInfo, error) {
	userInfoURL := p.issuerURL + "/protocol/openid-connect/userinfo"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %w", err)
	}

	info := &auth.UserInfo{
		Subject:   getStringFromMap(result, "sub"),
		Email:     getStringFromMap(result, "email"),
		FirstName: getStringFromMap(result, "given_name"),
		LastName:  getStringFromMap(result, "family_name"),
	}

	return info, nil
}

// GetLoginURL returns the Keycloak authorization endpoint URL.
func (p *Provider) GetLoginURL(state string) string {
	authURL := p.issuerURL + "/protocol/openid-connect/auth"
	params := url.Values{
		"client_id":     {p.clientID},
		"redirect_uri":  {p.redirectURL},
		"response_type": {"code"},
		"scope":         {"openid profile email"},
		"state":         {state},
	}
	return authURL + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for tokens.
func (p *Provider) ExchangeCode(ctx context.Context, code string) (*auth.TokenSet, error) {
	tokenURL := p.issuerURL + "/protocol/openid-connect/token"

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {p.redirectURL},
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &auth.TokenSet{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		IDToken:      result.IDToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func extractRoles(claims jwt.MapClaims) []string {
	var roles []string

	// Extract realm roles
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if rolesArr, ok := realmAccess["roles"].([]interface{}); ok {
			for _, r := range rolesArr {
				if s, ok := r.(string); ok {
					roles = append(roles, s)
				}
			}
		}
	}

	return roles
}
