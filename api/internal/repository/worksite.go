package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

type WorksiteRepository interface {
	List(ctx context.Context, companyID string, limit, offset int) ([]model.Worksite, error)
	GetByID(ctx context.Context, id string) (*model.Worksite, error)
	Create(ctx context.Context, worksite *model.Worksite) error
	Update(ctx context.Context, worksite *model.Worksite) error
	Delete(ctx context.Context, id string) error
}

type worksiteRepo struct {
	db *sql.DB
}

func NewWorksiteRepository(db *sql.DB) WorksiteRepository {
	return &worksiteRepo{db: db}
}

func (r *worksiteRepo) List(ctx context.Context, companyID string, limit, offset int) ([]model.Worksite, error) {
	query := `SELECT id, company_id, name, address, latitude, longitude, created_at, updated_at
		FROM worksites WHERE company_id = $1 ORDER BY name LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list worksites: %w", err)
	}
	defer rows.Close()

	var worksites []model.Worksite
	for rows.Next() {
		var w model.Worksite
		if err := rows.Scan(&w.ID, &w.CompanyID, &w.Name, &w.Address, &w.Latitude, &w.Longitude, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan worksite: %w", err)
		}
		worksites = append(worksites, w)
	}
	return worksites, rows.Err()
}

func (r *worksiteRepo) GetByID(ctx context.Context, id string) (*model.Worksite, error) {
	var w model.Worksite
	err := r.db.QueryRowContext(ctx,
		`SELECT id, company_id, name, address, latitude, longitude, created_at, updated_at
		FROM worksites WHERE id = $1`, id).
		Scan(&w.ID, &w.CompanyID, &w.Name, &w.Address, &w.Latitude, &w.Longitude, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worksite: %w", err)
	}
	return &w, nil
}

func (r *worksiteRepo) Create(ctx context.Context, worksite *model.Worksite) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO worksites (company_id, name, address, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		worksite.CompanyID, worksite.Name, worksite.Address, worksite.Latitude, worksite.Longitude).
		Scan(&worksite.ID, &worksite.CreatedAt, &worksite.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create worksite: %w", err)
	}
	return nil
}

func (r *worksiteRepo) Update(ctx context.Context, worksite *model.Worksite) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE worksites SET name = $1, address = $2, latitude = $3, longitude = $4, updated_at = NOW()
		WHERE id = $5`,
		worksite.Name, worksite.Address, worksite.Latitude, worksite.Longitude, worksite.ID)
	if err != nil {
		return fmt.Errorf("failed to update worksite: %w", err)
	}
	return nil
}

func (r *worksiteRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM worksites WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete worksite: %w", err)
	}
	return nil
}
