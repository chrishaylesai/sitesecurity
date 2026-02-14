package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// mockCompanyRepo is a test double for repository.CompanyRepository.
type mockCompanyRepo struct {
	companies []model.Company
	err       error
}

func (m *mockCompanyRepo) List(ctx context.Context, limit, offset int) ([]model.Company, error) {
	if m.err != nil {
		return nil, m.err
	}
	end := offset + limit
	if end > len(m.companies) {
		end = len(m.companies)
	}
	if offset >= len(m.companies) {
		return nil, nil
	}
	return m.companies[offset:end], nil
}

func (m *mockCompanyRepo) GetByID(ctx context.Context, id string) (*model.Company, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, c := range m.companies {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (m *mockCompanyRepo) Create(ctx context.Context, company *model.Company) error {
	if m.err != nil {
		return m.err
	}
	company.ID = "new-id"
	m.companies = append(m.companies, *company)
	return nil
}

func (m *mockCompanyRepo) Update(ctx context.Context, company *model.Company) error {
	return m.err
}

func (m *mockCompanyRepo) Delete(ctx context.Context, id string) error {
	return m.err
}

func TestCompanyService_List(t *testing.T) {
	repo := &mockCompanyRepo{
		companies: []model.Company{
			{ID: "1", Name: "Alpha Security"},
			{ID: "2", Name: "Beta Guard"},
		},
	}
	svc := service.NewCompanyService(repo)

	companies, err := svc.List(context.Background(), 1, 25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(companies) != 2 {
		t.Errorf("expected 2 companies, got %d", len(companies))
	}
}

func TestCompanyService_List_DefaultPagination(t *testing.T) {
	repo := &mockCompanyRepo{companies: []model.Company{}}
	svc := service.NewCompanyService(repo)

	_, err := svc.List(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompanyService_GetByID_Found(t *testing.T) {
	repo := &mockCompanyRepo{
		companies: []model.Company{{ID: "abc", Name: "Test Corp"}},
	}
	svc := service.NewCompanyService(repo)

	company, err := svc.GetByID(context.Background(), "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if company.Name != "Test Corp" {
		t.Errorf("expected 'Test Corp', got '%s'", company.Name)
	}
}

func TestCompanyService_GetByID_NotFound(t *testing.T) {
	repo := &mockCompanyRepo{companies: []model.Company{}}
	svc := service.NewCompanyService(repo)

	_, err := svc.GetByID(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing company")
	}
}

func TestCompanyService_Create_Valid(t *testing.T) {
	repo := &mockCompanyRepo{}
	svc := service.NewCompanyService(repo)

	company := &model.Company{Name: "New Security Co"}
	err := svc.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if company.ID != "new-id" {
		t.Errorf("expected ID to be set, got '%s'", company.ID)
	}
}

func TestCompanyService_Create_EmptyName(t *testing.T) {
	repo := &mockCompanyRepo{}
	svc := service.NewCompanyService(repo)

	company := &model.Company{Name: ""}
	err := svc.Create(context.Background(), company)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestCompanyService_Delete_NotFound(t *testing.T) {
	repo := &mockCompanyRepo{companies: []model.Company{}}
	svc := service.NewCompanyService(repo)

	err := svc.Delete(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing company")
	}
}

func TestCompanyService_List_RepoError(t *testing.T) {
	repo := &mockCompanyRepo{err: fmt.Errorf("database error")}
	svc := service.NewCompanyService(repo)

	_, err := svc.List(context.Background(), 1, 25)
	if err == nil {
		t.Error("expected error from repo")
	}
}
