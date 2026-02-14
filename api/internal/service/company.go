package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

// CompanyService handles business logic for companies.
type CompanyService struct {
	repo repository.CompanyRepository
}

// NewCompanyService creates a new CompanyService.
func NewCompanyService(repo repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) List(ctx context.Context, page, perPage int) ([]model.Company, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.repo.List(ctx, perPage, offset)
}

func (s *CompanyService) GetByID(ctx context.Context, id string) (*model.Company, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}
	return company, nil
}

func (s *CompanyService) Create(ctx context.Context, company *model.Company) error {
	if company.Name == "" {
		return fmt.Errorf("company name is required")
	}
	return s.repo.Create(ctx, company)
}

func (s *CompanyService) Update(ctx context.Context, company *model.Company) error {
	if company.Name == "" {
		return fmt.Errorf("company name is required")
	}
	existing, err := s.repo.GetByID(ctx, company.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("company not found")
	}
	return s.repo.Update(ctx, company)
}

func (s *CompanyService) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("company not found")
	}
	return s.repo.Delete(ctx, id)
}
