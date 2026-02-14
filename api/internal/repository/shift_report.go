package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

// ShiftReportTemplateRepository defines the interface for shift report template data access.
type ShiftReportTemplateRepository interface {
	ListByCompany(ctx context.Context, companyID string, limit, offset int) ([]model.ShiftReportTemplate, error)
	GetByID(ctx context.Context, id string) (*model.ShiftReportTemplate, error)
	Create(ctx context.Context, template *model.ShiftReportTemplate) error
	Update(ctx context.Context, template *model.ShiftReportTemplate) error
	Delete(ctx context.Context, id string) error
}

type shiftReportTemplateRepo struct {
	db *sql.DB
}

// NewShiftReportTemplateRepository creates a new ShiftReportTemplateRepository.
func NewShiftReportTemplateRepository(db *sql.DB) ShiftReportTemplateRepository {
	return &shiftReportTemplateRepo{db: db}
}

func (r *shiftReportTemplateRepo) ListByCompany(ctx context.Context, companyID string, limit, offset int) ([]model.ShiftReportTemplate, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, company_id, name, fields, created_at, updated_at
		 FROM shift_report_templates WHERE company_id = $1 ORDER BY name LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shift report templates: %w", err)
	}
	defer rows.Close()

	var templates []model.ShiftReportTemplate
	for rows.Next() {
		var t model.ShiftReportTemplate
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Fields, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift report template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *shiftReportTemplateRepo) GetByID(ctx context.Context, id string) (*model.ShiftReportTemplate, error) {
	var t model.ShiftReportTemplate
	err := r.db.QueryRowContext(ctx,
		`SELECT id, company_id, name, fields, created_at, updated_at
		 FROM shift_report_templates WHERE id = $1`, id).
		Scan(&t.ID, &t.CompanyID, &t.Name, &t.Fields, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shift report template: %w", err)
	}
	return &t, nil
}

func (r *shiftReportTemplateRepo) Create(ctx context.Context, template *model.ShiftReportTemplate) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO shift_report_templates (company_id, name, fields)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		template.CompanyID, template.Name, template.Fields).
		Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create shift report template: %w", err)
	}
	return nil
}

func (r *shiftReportTemplateRepo) Update(ctx context.Context, template *model.ShiftReportTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE shift_report_templates SET name = $1, fields = $2, updated_at = NOW()
		 WHERE id = $3`,
		template.Name, template.Fields, template.ID)
	if err != nil {
		return fmt.Errorf("failed to update shift report template: %w", err)
	}
	return nil
}

func (r *shiftReportTemplateRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shift_report_templates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete shift report template: %w", err)
	}
	return nil
}

// ShiftReportRepository defines the interface for shift report data access.
type ShiftReportRepository interface {
	ListByShift(ctx context.Context, shiftID string, limit, offset int) ([]model.ShiftReport, error)
	ListByWorker(ctx context.Context, workerID string, limit, offset int) ([]model.ShiftReport, error)
	GetByID(ctx context.Context, id string) (*model.ShiftReport, error)
	Create(ctx context.Context, report *model.ShiftReport) error
}

type shiftReportRepo struct {
	db *sql.DB
}

// NewShiftReportRepository creates a new ShiftReportRepository.
func NewShiftReportRepository(db *sql.DB) ShiftReportRepository {
	return &shiftReportRepo{db: db}
}

func (r *shiftReportRepo) ListByShift(ctx context.Context, shiftID string, limit, offset int) ([]model.ShiftReport, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shift_id, worker_id, template_id, data, submitted_at
		 FROM shift_reports WHERE shift_id = $1 ORDER BY submitted_at DESC LIMIT $2 OFFSET $3`,
		shiftID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shift reports by shift: %w", err)
	}
	defer rows.Close()

	var reports []model.ShiftReport
	for rows.Next() {
		var sr model.ShiftReport
		if err := rows.Scan(&sr.ID, &sr.ShiftID, &sr.WorkerID, &sr.TemplateID, &sr.Data, &sr.SubmittedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift report: %w", err)
		}
		reports = append(reports, sr)
	}
	return reports, rows.Err()
}

func (r *shiftReportRepo) ListByWorker(ctx context.Context, workerID string, limit, offset int) ([]model.ShiftReport, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shift_id, worker_id, template_id, data, submitted_at
		 FROM shift_reports WHERE worker_id = $1 ORDER BY submitted_at DESC LIMIT $2 OFFSET $3`,
		workerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shift reports by worker: %w", err)
	}
	defer rows.Close()

	var reports []model.ShiftReport
	for rows.Next() {
		var sr model.ShiftReport
		if err := rows.Scan(&sr.ID, &sr.ShiftID, &sr.WorkerID, &sr.TemplateID, &sr.Data, &sr.SubmittedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift report: %w", err)
		}
		reports = append(reports, sr)
	}
	return reports, rows.Err()
}

func (r *shiftReportRepo) GetByID(ctx context.Context, id string) (*model.ShiftReport, error) {
	var sr model.ShiftReport
	err := r.db.QueryRowContext(ctx,
		`SELECT id, shift_id, worker_id, template_id, data, submitted_at
		 FROM shift_reports WHERE id = $1`, id).
		Scan(&sr.ID, &sr.ShiftID, &sr.WorkerID, &sr.TemplateID, &sr.Data, &sr.SubmittedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shift report: %w", err)
	}
	return &sr, nil
}

func (r *shiftReportRepo) Create(ctx context.Context, report *model.ShiftReport) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO shift_reports (shift_id, worker_id, template_id, data)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, submitted_at`,
		report.ShiftID, report.WorkerID, report.TemplateID, report.Data).
		Scan(&report.ID, &report.SubmittedAt)
	if err != nil {
		return fmt.Errorf("failed to create shift report: %w", err)
	}
	return nil
}
