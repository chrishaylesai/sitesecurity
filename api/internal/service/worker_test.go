package service_test

import (
	"context"
	"testing"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

type mockWorkerRepo struct {
	workers []model.Worker
	err     error
}

func (m *mockWorkerRepo) List(ctx context.Context, limit, offset int) ([]model.Worker, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.workers, nil
}

func (m *mockWorkerRepo) GetByID(ctx context.Context, id string) (*model.Worker, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, w := range m.workers {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, nil
}

func (m *mockWorkerRepo) GetByAuthSubject(ctx context.Context, authSubject string) (*model.Worker, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, w := range m.workers {
		if w.AuthSubject == authSubject {
			return &w, nil
		}
	}
	return nil, nil
}

func (m *mockWorkerRepo) Create(ctx context.Context, worker *model.Worker) error {
	if m.err != nil {
		return m.err
	}
	worker.ID = "new-worker-id"
	m.workers = append(m.workers, *worker)
	return nil
}

func (m *mockWorkerRepo) Update(ctx context.Context, worker *model.Worker) error { return m.err }
func (m *mockWorkerRepo) Delete(ctx context.Context, id string) error             { return m.err }

type mockCertRepo struct {
	certs []model.Certificate
	err   error
}

func (m *mockCertRepo) ListByWorker(ctx context.Context, workerID string) ([]model.Certificate, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.Certificate
	for _, c := range m.certs {
		if c.WorkerID == workerID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCertRepo) GetByID(ctx context.Context, id string) (*model.Certificate, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, c := range m.certs {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (m *mockCertRepo) Create(ctx context.Context, cert *model.Certificate) error {
	if m.err != nil {
		return m.err
	}
	cert.ID = "new-cert-id"
	return nil
}

func (m *mockCertRepo) Update(ctx context.Context, cert *model.Certificate) error { return m.err }
func (m *mockCertRepo) Delete(ctx context.Context, id string) error               { return m.err }

type mockWCRepo struct {
	memberships []model.WorkerCompany
	err         error
}

func (m *mockWCRepo) ListByWorker(ctx context.Context, workerID string) ([]model.WorkerCompany, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.WorkerCompany
	for _, wc := range m.memberships {
		if wc.WorkerID == workerID {
			result = append(result, wc)
		}
	}
	return result, nil
}

func (m *mockWCRepo) ListByCompany(ctx context.Context, companyID string) ([]model.WorkerCompany, error) {
	return nil, m.err
}

func (m *mockWCRepo) Get(ctx context.Context, workerID, companyID string) (*model.WorkerCompany, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, wc := range m.memberships {
		if wc.WorkerID == workerID && wc.CompanyID == companyID {
			return &wc, nil
		}
	}
	return nil, nil
}

func (m *mockWCRepo) Create(ctx context.Context, wc *model.WorkerCompany) error { return m.err }
func (m *mockWCRepo) UpdateRole(ctx context.Context, workerID, companyID string, role model.WorkerRole) error {
	return m.err
}
func (m *mockWCRepo) UpdateStatus(ctx context.Context, workerID, companyID string, status model.MembershipStatus) error {
	return m.err
}
func (m *mockWCRepo) Delete(ctx context.Context, workerID, companyID string) error { return m.err }

func TestWorkerService_Create_Valid(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	worker := &model.Worker{
		AuthSubject: "sub-123",
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john@example.com",
	}
	err := svc.Create(context.Background(), worker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if worker.ID != "new-worker-id" {
		t.Errorf("expected ID to be set")
	}
}

func TestWorkerService_Create_MissingName(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	err := svc.Create(context.Background(), &model.Worker{
		AuthSubject: "sub-123",
		Email:       "test@example.com",
	})
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestWorkerService_Create_MissingEmail(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	err := svc.Create(context.Background(), &model.Worker{
		AuthSubject: "sub-123",
		FirstName:   "John",
		LastName:    "Smith",
	})
	if err == nil {
		t.Error("expected error for missing email")
	}
}

func TestWorkerService_Create_MissingAuthSubject(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	err := svc.Create(context.Background(), &model.Worker{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
	})
	if err == nil {
		t.Error("expected error for missing auth subject")
	}
}

func TestWorkerService_GetByID_NotFound(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	_, err := svc.GetByID(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing worker")
	}
}

func TestWorkerService_CreateCertificate_Valid(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	cert := &model.Certificate{WorkerID: "w1", Name: "SIA Door Supervisor"}
	err := svc.CreateCertificate(context.Background(), cert)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.ID != "new-cert-id" {
		t.Errorf("expected ID to be set")
	}
}

func TestWorkerService_CreateCertificate_MissingName(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	err := svc.CreateCertificate(context.Background(), &model.Certificate{WorkerID: "w1"})
	if err == nil {
		t.Error("expected error for missing cert name")
	}
}

func TestWorkerService_AddMembership_Duplicate(t *testing.T) {
	wcRepo := &mockWCRepo{
		memberships: []model.WorkerCompany{
			{WorkerID: "w1", CompanyID: "c1", Role: model.RoleWorker, Status: model.MembershipActive},
		},
	}
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, wcRepo)
	err := svc.AddMembership(context.Background(), &model.WorkerCompany{WorkerID: "w1", CompanyID: "c1"})
	if err == nil {
		t.Error("expected error for duplicate membership")
	}
}

func TestWorkerService_RemoveMembership_NotFound(t *testing.T) {
	svc := service.NewWorkerService(&mockWorkerRepo{}, &mockCertRepo{}, &mockWCRepo{})
	err := svc.RemoveMembership(context.Background(), "w1", "c1")
	if err == nil {
		t.Error("expected error for missing membership")
	}
}
