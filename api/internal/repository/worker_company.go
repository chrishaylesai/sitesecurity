package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

type WorkerCompanyRepository interface {
	ListByWorker(ctx context.Context, workerID string) ([]model.WorkerCompany, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.WorkerCompany, error)
	Get(ctx context.Context, workerID, companyID string) (*model.WorkerCompany, error)
	Create(ctx context.Context, wc *model.WorkerCompany) error
	UpdateRole(ctx context.Context, workerID, companyID string, role model.WorkerRole) error
	UpdateStatus(ctx context.Context, workerID, companyID string, status model.MembershipStatus) error
	Delete(ctx context.Context, workerID, companyID string) error
}

type workerCompanyRepo struct {
	db *sql.DB
}

func NewWorkerCompanyRepository(db *sql.DB) WorkerCompanyRepository {
	return &workerCompanyRepo{db: db}
}

func (r *workerCompanyRepo) ListByWorker(ctx context.Context, workerID string) ([]model.WorkerCompany, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT worker_id, company_id, role, status, joined_at
		FROM worker_companies WHERE worker_id = $1`, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list worker companies: %w", err)
	}
	defer rows.Close()

	var memberships []model.WorkerCompany
	for rows.Next() {
		var wc model.WorkerCompany
		if err := rows.Scan(&wc.WorkerID, &wc.CompanyID, &wc.Role, &wc.Status, &wc.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan worker company: %w", err)
		}
		memberships = append(memberships, wc)
	}
	return memberships, rows.Err()
}

func (r *workerCompanyRepo) ListByCompany(ctx context.Context, companyID string) ([]model.WorkerCompany, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT worker_id, company_id, role, status, joined_at
		FROM worker_companies WHERE company_id = $1`, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list company workers: %w", err)
	}
	defer rows.Close()

	var memberships []model.WorkerCompany
	for rows.Next() {
		var wc model.WorkerCompany
		if err := rows.Scan(&wc.WorkerID, &wc.CompanyID, &wc.Role, &wc.Status, &wc.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan worker company: %w", err)
		}
		memberships = append(memberships, wc)
	}
	return memberships, rows.Err()
}

func (r *workerCompanyRepo) Get(ctx context.Context, workerID, companyID string) (*model.WorkerCompany, error) {
	var wc model.WorkerCompany
	err := r.db.QueryRowContext(ctx,
		`SELECT worker_id, company_id, role, status, joined_at
		FROM worker_companies WHERE worker_id = $1 AND company_id = $2`, workerID, companyID).
		Scan(&wc.WorkerID, &wc.CompanyID, &wc.Role, &wc.Status, &wc.JoinedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worker company: %w", err)
	}
	return &wc, nil
}

func (r *workerCompanyRepo) Create(ctx context.Context, wc *model.WorkerCompany) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO worker_companies (worker_id, company_id, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING joined_at`,
		wc.WorkerID, wc.CompanyID, wc.Role, wc.Status).
		Scan(&wc.JoinedAt)
	if err != nil {
		return fmt.Errorf("failed to create worker company: %w", err)
	}
	return nil
}

func (r *workerCompanyRepo) UpdateRole(ctx context.Context, workerID, companyID string, role model.WorkerRole) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE worker_companies SET role = $1 WHERE worker_id = $2 AND company_id = $3`,
		role, workerID, companyID)
	if err != nil {
		return fmt.Errorf("failed to update worker role: %w", err)
	}
	return nil
}

func (r *workerCompanyRepo) UpdateStatus(ctx context.Context, workerID, companyID string, status model.MembershipStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE worker_companies SET status = $1 WHERE worker_id = $2 AND company_id = $3`,
		status, workerID, companyID)
	if err != nil {
		return fmt.Errorf("failed to update membership status: %w", err)
	}
	return nil
}

func (r *workerCompanyRepo) Delete(ctx context.Context, workerID, companyID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM worker_companies WHERE worker_id = $1 AND company_id = $2`,
		workerID, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete worker company: %w", err)
	}
	return nil
}
