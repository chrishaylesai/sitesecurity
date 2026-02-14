package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

// ShiftReportService handles business logic for shift reports and templates.
type ShiftReportService struct {
	templateRepo repository.ShiftReportTemplateRepository
	reportRepo   repository.ShiftReportRepository
}

// NewShiftReportService creates a new ShiftReportService.
func NewShiftReportService(templateRepo repository.ShiftReportTemplateRepository, reportRepo repository.ShiftReportRepository) *ShiftReportService {
	return &ShiftReportService{templateRepo: templateRepo, reportRepo: reportRepo}
}

// ListTemplates returns templates for a company with pagination.
func (s *ShiftReportService) ListTemplates(ctx context.Context, companyID string, page, perPage int) ([]model.ShiftReportTemplate, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.templateRepo.ListByCompany(ctx, companyID, perPage, offset)
}

// GetTemplateByID returns a single template by ID.
func (s *ShiftReportService) GetTemplateByID(ctx context.Context, id string) (*model.ShiftReportTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("shift report template not found")
	}
	return template, nil
}

// CreateTemplate creates a new shift report template.
func (s *ShiftReportService) CreateTemplate(ctx context.Context, template *model.ShiftReportTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}
	if template.CompanyID == "" {
		return fmt.Errorf("company_id is required")
	}
	return s.templateRepo.Create(ctx, template)
}

// UpdateTemplate updates an existing shift report template.
func (s *ShiftReportService) UpdateTemplate(ctx context.Context, template *model.ShiftReportTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}
	existing, err := s.templateRepo.GetByID(ctx, template.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("shift report template not found")
	}
	return s.templateRepo.Update(ctx, template)
}

// DeleteTemplate deletes a shift report template.
func (s *ShiftReportService) DeleteTemplate(ctx context.Context, id string) error {
	existing, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("shift report template not found")
	}
	return s.templateRepo.Delete(ctx, id)
}

// ListReportsByShift returns reports for a shift with pagination.
func (s *ShiftReportService) ListReportsByShift(ctx context.Context, shiftID string, page, perPage int) ([]model.ShiftReport, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.reportRepo.ListByShift(ctx, shiftID, perPage, offset)
}

// ListReportsByWorker returns reports for a worker with pagination.
func (s *ShiftReportService) ListReportsByWorker(ctx context.Context, workerID string, page, perPage int) ([]model.ShiftReport, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.reportRepo.ListByWorker(ctx, workerID, perPage, offset)
}

// GetReportByID returns a single report by ID.
func (s *ShiftReportService) GetReportByID(ctx context.Context, id string) (*model.ShiftReport, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, fmt.Errorf("shift report not found")
	}
	return report, nil
}

// CreateReport creates a new shift report.
func (s *ShiftReportService) CreateReport(ctx context.Context, report *model.ShiftReport) error {
	if report.ShiftID == "" {
		return fmt.Errorf("shift_id is required")
	}
	if report.WorkerID == "" {
		return fmt.Errorf("worker_id is required")
	}
	return s.reportRepo.Create(ctx, report)
}
