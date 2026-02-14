package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

type WorkerService struct {
	workerRepo repository.WorkerRepository
	certRepo   repository.CertificateRepository
	wcRepo     repository.WorkerCompanyRepository
}

func NewWorkerService(
	workerRepo repository.WorkerRepository,
	certRepo repository.CertificateRepository,
	wcRepo repository.WorkerCompanyRepository,
) *WorkerService {
	return &WorkerService{
		workerRepo: workerRepo,
		certRepo:   certRepo,
		wcRepo:     wcRepo,
	}
}

func (s *WorkerService) List(ctx context.Context, page, perPage int) ([]model.Worker, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.workerRepo.List(ctx, perPage, offset)
}

func (s *WorkerService) GetByID(ctx context.Context, id string) (*model.Worker, error) {
	worker, err := s.workerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if worker == nil {
		return nil, fmt.Errorf("worker not found")
	}
	return worker, nil
}

func (s *WorkerService) GetByAuthSubject(ctx context.Context, authSubject string) (*model.Worker, error) {
	return s.workerRepo.GetByAuthSubject(ctx, authSubject)
}

func (s *WorkerService) Create(ctx context.Context, worker *model.Worker) error {
	if worker.FirstName == "" || worker.LastName == "" {
		return fmt.Errorf("first name and last name are required")
	}
	if worker.Email == "" {
		return fmt.Errorf("email is required")
	}
	if worker.AuthSubject == "" {
		return fmt.Errorf("auth subject is required")
	}
	return s.workerRepo.Create(ctx, worker)
}

func (s *WorkerService) Update(ctx context.Context, worker *model.Worker) error {
	existing, err := s.workerRepo.GetByID(ctx, worker.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("worker not found")
	}
	return s.workerRepo.Update(ctx, worker)
}

// Certificates

func (s *WorkerService) ListCertificates(ctx context.Context, workerID string) ([]model.Certificate, error) {
	return s.certRepo.ListByWorker(ctx, workerID)
}

func (s *WorkerService) GetCertificate(ctx context.Context, id string) (*model.Certificate, error) {
	cert, err := s.certRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, fmt.Errorf("certificate not found")
	}
	return cert, nil
}

func (s *WorkerService) CreateCertificate(ctx context.Context, cert *model.Certificate) error {
	if cert.Name == "" {
		return fmt.Errorf("certificate name is required")
	}
	if cert.WorkerID == "" {
		return fmt.Errorf("worker ID is required")
	}
	return s.certRepo.Create(ctx, cert)
}

func (s *WorkerService) UpdateCertificate(ctx context.Context, cert *model.Certificate) error {
	existing, err := s.certRepo.GetByID(ctx, cert.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("certificate not found")
	}
	return s.certRepo.Update(ctx, cert)
}

func (s *WorkerService) DeleteCertificate(ctx context.Context, id string) error {
	existing, err := s.certRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("certificate not found")
	}
	return s.certRepo.Delete(ctx, id)
}

// Memberships

func (s *WorkerService) ListMemberships(ctx context.Context, workerID string) ([]model.WorkerCompany, error) {
	return s.wcRepo.ListByWorker(ctx, workerID)
}

func (s *WorkerService) ListCompanyMembers(ctx context.Context, companyID string) ([]model.WorkerCompany, error) {
	return s.wcRepo.ListByCompany(ctx, companyID)
}

func (s *WorkerService) AddMembership(ctx context.Context, wc *model.WorkerCompany) error {
	if wc.WorkerID == "" || wc.CompanyID == "" {
		return fmt.Errorf("worker ID and company ID are required")
	}
	existing, err := s.wcRepo.Get(ctx, wc.WorkerID, wc.CompanyID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("membership already exists")
	}
	if wc.Status == "" {
		wc.Status = model.MembershipActive
	}
	if wc.Role == "" {
		wc.Role = model.RoleWorker
	}
	return s.wcRepo.Create(ctx, wc)
}

func (s *WorkerService) UpdateMembershipRole(ctx context.Context, workerID, companyID string, role model.WorkerRole) error {
	existing, err := s.wcRepo.Get(ctx, workerID, companyID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("membership not found")
	}
	return s.wcRepo.UpdateRole(ctx, workerID, companyID, role)
}

func (s *WorkerService) RemoveMembership(ctx context.Context, workerID, companyID string) error {
	existing, err := s.wcRepo.Get(ctx, workerID, companyID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("membership not found")
	}
	return s.wcRepo.Delete(ctx, workerID, companyID)
}
