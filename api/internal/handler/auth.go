package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/auth"
)

// AuthHandler handles authentication routes.
type AuthHandler struct {
	provider auth.Provider
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(provider auth.Provider) *AuthHandler {
	return &AuthHandler{provider: provider}
}

// Routes returns the auth routes.
func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/login", h.Login)
	r.Get("/callback", h.Callback)
	return r
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	loginURL := h.provider.GetLoginURL(state)
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		Error(w, http.StatusBadRequest, "missing authorization code")
		return
	}

	tokenSet, err := h.provider.ExchangeCode(r.Context(), code)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to exchange authorization code")
		return
	}

	JSON(w, http.StatusOK, tokenSet)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
