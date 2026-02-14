package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

type CertificateRepository interface {
	ListByWorker(ctx context.Context, workerID string) ([]model.Certificate, error)
	GetByID(ctx context.Context, id string) (*model.Certificate, error)
	Create(ctx context.Context, cert *model.Certificate) error
	Update(ctx context.Context, cert *model.Certificate) error
	Delete(ctx context.Context, id string) error
}

type certificateRepo struct {
	db *sql.DB
}

func NewCertificateRepository(db *sql.DB) CertificateRepository {
	return &certificateRepo{db: db}
}

func (r *certificateRepo) ListByWorker(ctx context.Context, workerID string) ([]model.Certificate, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worker_id, name, issuing_body, certificate_number, issued_date, expiry_date, created_at, updated_at
		FROM certificates WHERE worker_id = $1 ORDER BY expiry_date DESC`, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list certificates: %w", err)
	}
	defer rows.Close()

	var certs []model.Certificate
	for rows.Next() {
		var c model.Certificate
		if err := rows.Scan(&c.ID, &c.WorkerID, &c.Name, &c.IssuingBody, &c.CertificateNumber, &c.IssuedDate, &c.ExpiryDate, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan certificate: %w", err)
		}
		certs = append(certs, c)
	}
	return certs, rows.Err()
}

func (r *certificateRepo) GetByID(ctx context.Context, id string) (*model.Certificate, error) {
	var c model.Certificate
	err := r.db.QueryRowContext(ctx,
		`SELECT id, worker_id, name, issuing_body, certificate_number, issued_date, expiry_date, created_at, updated_at
		FROM certificates WHERE id = $1`, id).
		Scan(&c.ID, &c.WorkerID, &c.Name, &c.IssuingBody, &c.CertificateNumber, &c.IssuedDate, &c.ExpiryDate, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}
	return &c, nil
}

func (r *certificateRepo) Create(ctx context.Context, cert *model.Certificate) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO certificates (worker_id, name, issuing_body, certificate_number, issued_date, expiry_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		cert.WorkerID, cert.Name, cert.IssuingBody, cert.CertificateNumber, cert.IssuedDate, cert.ExpiryDate).
		Scan(&cert.ID, &cert.CreatedAt, &cert.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}
	return nil
}

func (r *certificateRepo) Update(ctx context.Context, cert *model.Certificate) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE certificates SET name = $1, issuing_body = $2, certificate_number = $3, issued_date = $4, expiry_date = $5, updated_at = NOW()
		WHERE id = $6`,
		cert.Name, cert.IssuingBody, cert.CertificateNumber, cert.IssuedDate, cert.ExpiryDate, cert.ID)
	if err != nil {
		return fmt.Errorf("failed to update certificate: %w", err)
	}
	return nil
}

func (r *certificateRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM certificates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}
	return nil
}
