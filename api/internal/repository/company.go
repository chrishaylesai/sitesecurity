package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

// CompanyRepository defines the interface for company data access.
type CompanyRepository interface {
	List(ctx context.Context, limit, offset int) ([]model.Company, error)
	GetByID(ctx context.Context, id string) (*model.Company, error)
	Create(ctx context.Context, company *model.Company) error
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, id string) error
}

type companyRepo struct {
	db *sql.DB
}

// NewCompanyRepository creates a new CompanyRepository.
func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &companyRepo{db: db}
}

func (r *companyRepo) List(ctx context.Context, limit, offset int) ([]model.Company, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, phone, email, created_at, updated_at
		 FROM companies ORDER BY name LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}
	defer rows.Close()

	var companies []model.Company
	for rows.Next() {
		var c model.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

func (r *companyRepo) GetByID(ctx context.Context, id string) (*model.Company, error) {
	var c model.Company
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, address, phone, email, created_at, updated_at
		 FROM companies WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.Address, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}
	return &c, nil
}

func (r *companyRepo) Create(ctx context.Context, company *model.Company) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO companies (name, address, phone, email)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, updated_at`,
		company.Name, company.Address, company.Phone, company.Email).
		Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

func (r *companyRepo) Update(ctx context.Context, company *model.Company) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE companies SET name = $1, address = $2, phone = $3, email = $4, updated_at = NOW()
		 WHERE id = $5`,
		company.Name, company.Address, company.Phone, company.Email, company.ID)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}

func (r *companyRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM companies WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}
