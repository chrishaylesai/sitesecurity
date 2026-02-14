package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

type WorksiteService struct {
	repo repository.WorksiteRepository
}

func NewWorksiteService(repo repository.WorksiteRepository) *WorksiteService {
	return &WorksiteService{repo: repo}
}

func (s *WorksiteService) List(ctx context.Context, companyID string, page, perPage int) ([]model.Worksite, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company ID is required")
	}
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.repo.List(ctx, companyID, perPage, offset)
}

func (s *WorksiteService) GetByID(ctx context.Context, id string) (*model.Worksite, error) {
	worksite, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if worksite == nil {
		return nil, fmt.Errorf("worksite not found")
	}
	return worksite, nil
}

func (s *WorksiteService) Create(ctx context.Context, worksite *model.Worksite) error {
	if worksite.Name == "" {
		return fmt.Errorf("worksite name is required")
	}
	if worksite.CompanyID == "" {
		return fmt.Errorf("company ID is required")
	}
	return s.repo.Create(ctx, worksite)
}

func (s *WorksiteService) Update(ctx context.Context, worksite *model.Worksite) error {
	if worksite.Name == "" {
		return fmt.Errorf("worksite name is required")
	}
	existing, err := s.repo.GetByID(ctx, worksite.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("worksite not found")
	}
	return s.repo.Update(ctx, worksite)
}

func (s *WorksiteService) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("worksite not found")
	}
	return s.repo.Delete(ctx, id)
}
