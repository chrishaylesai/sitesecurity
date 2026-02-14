package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/handler"
	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var result map[string]string
	json.NewDecoder(rr.Body).Decode(&result)
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", result["status"])
	}
}

func TestCompanyHandler_Create(t *testing.T) {
	repo := &inMemoryCompanyRepo{companies: make(map[string]model.Company)}
	svc := service.NewCompanyService(repo)
	h := handler.NewCompanyHandler(svc)

	body := `{"name":"Test Security Ltd"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Mount("/", h.Routes())
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var company model.Company
	json.NewDecoder(rr.Body).Decode(&company)
	if company.Name != "Test Security Ltd" {
		t.Errorf("expected name 'Test Security Ltd', got '%s'", company.Name)
	}
	if company.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestCompanyHandler_Create_InvalidBody(t *testing.T) {
	repo := &inMemoryCompanyRepo{companies: make(map[string]model.Company)}
	svc := service.NewCompanyService(repo)
	h := handler.NewCompanyHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Mount("/", h.Routes())
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCompanyHandler_List(t *testing.T) {
	repo := &inMemoryCompanyRepo{companies: make(map[string]model.Company)}
	repo.companies["1"] = model.Company{ID: "1", Name: "Alpha"}
	repo.companies["2"] = model.Company{ID: "2", Name: "Beta"}

	svc := service.NewCompanyService(repo)
	h := handler.NewCompanyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Mount("/", h.Routes())
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var companies []model.Company
	json.NewDecoder(rr.Body).Decode(&companies)
	if len(companies) != 2 {
		t.Errorf("expected 2 companies, got %d", len(companies))
	}
}

func TestCompanyHandler_GetByID_NotFound(t *testing.T) {
	repo := &inMemoryCompanyRepo{companies: make(map[string]model.Company)}
	svc := service.NewCompanyService(repo)
	h := handler.NewCompanyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/missing-id", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Mount("/", h.Routes())
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

// inMemoryCompanyRepo implements repository.CompanyRepository for testing.
type inMemoryCompanyRepo struct {
	companies map[string]model.Company
	nextID    int
}

func (r *inMemoryCompanyRepo) List(ctx context.Context, limit, offset int) ([]model.Company, error) {
	var result []model.Company
	for _, c := range r.companies {
		result = append(result, c)
	}
	return result, nil
}

func (r *inMemoryCompanyRepo) GetByID(ctx context.Context, id string) (*model.Company, error) {
	c, ok := r.companies[id]
	if !ok {
		return nil, nil
	}
	return &c, nil
}

func (r *inMemoryCompanyRepo) Create(ctx context.Context, company *model.Company) error {
	r.nextID++
	company.ID = fmt.Sprintf("test-%d", r.nextID)
	r.companies[company.ID] = *company
	return nil
}

func (r *inMemoryCompanyRepo) Update(ctx context.Context, company *model.Company) error {
	r.companies[company.ID] = *company
	return nil
}

func (r *inMemoryCompanyRepo) Delete(ctx context.Context, id string) error {
	delete(r.companies, id)
	return nil
}
