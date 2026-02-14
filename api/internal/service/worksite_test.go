package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

type mockWorksiteRepo struct {
	worksites []model.Worksite
	err       error
}

func (m *mockWorksiteRepo) List(ctx context.Context, companyID string, limit, offset int) ([]model.Worksite, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.Worksite
	for _, w := range m.worksites {
		if w.CompanyID == companyID {
			result = append(result, w)
		}
	}
	return result, nil
}

func (m *mockWorksiteRepo) GetByID(ctx context.Context, id string) (*model.Worksite, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, w := range m.worksites {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, nil
}

func (m *mockWorksiteRepo) Create(ctx context.Context, worksite *model.Worksite) error {
	if m.err != nil {
		return m.err
	}
	worksite.ID = "new-ws-id"
	m.worksites = append(m.worksites, *worksite)
	return nil
}

func (m *mockWorksiteRepo) Update(ctx context.Context, worksite *model.Worksite) error { return m.err }
func (m *mockWorksiteRepo) Delete(ctx context.Context, id string) error                { return m.err }

func TestWorksiteService_List_RequiresCompanyID(t *testing.T) {
	svc := service.NewWorksiteService(&mockWorksiteRepo{})
	_, err := svc.List(context.Background(), "", 1, 25)
	if err == nil {
		t.Error("expected error for empty company ID")
	}
}

func TestWorksiteService_Create_Valid(t *testing.T) {
	repo := &mockWorksiteRepo{}
	svc := service.NewWorksiteService(repo)

	ws := &model.Worksite{Name: "Test Site", CompanyID: "company-1"}
	err := svc.Create(context.Background(), ws)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.ID != "new-ws-id" {
		t.Errorf("expected ID to be set")
	}
}

func TestWorksiteService_Create_EmptyName(t *testing.T) {
	svc := service.NewWorksiteService(&mockWorksiteRepo{})
	err := svc.Create(context.Background(), &model.Worksite{CompanyID: "c1"})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestWorksiteService_Create_EmptyCompanyID(t *testing.T) {
	svc := service.NewWorksiteService(&mockWorksiteRepo{})
	err := svc.Create(context.Background(), &model.Worksite{Name: "Site"})
	if err == nil {
		t.Error("expected error for empty company ID")
	}
}

func TestWorksiteService_GetByID_NotFound(t *testing.T) {
	svc := service.NewWorksiteService(&mockWorksiteRepo{})
	_, err := svc.GetByID(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing worksite")
	}
}

func TestWorksiteService_Delete_NotFound(t *testing.T) {
	svc := service.NewWorksiteService(&mockWorksiteRepo{})
	err := svc.Delete(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing worksite")
	}
}

func TestWorksiteService_List_RepoError(t *testing.T) {
	repo := &mockWorksiteRepo{err: fmt.Errorf("db error")}
	svc := service.NewWorksiteService(repo)
	_, err := svc.List(context.Background(), "c1", 1, 25)
	if err == nil {
		t.Error("expected error from repo")
	}
}
