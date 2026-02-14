package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

type WorkerRepository interface {
	List(ctx context.Context, limit, offset int) ([]model.Worker, error)
	GetByID(ctx context.Context, id string) (*model.Worker, error)
	GetByAuthSubject(ctx context.Context, authSubject string) (*model.Worker, error)
	Create(ctx context.Context, worker *model.Worker) error
	Update(ctx context.Context, worker *model.Worker) error
	Delete(ctx context.Context, id string) error
}

type workerRepo struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) WorkerRepository {
	return &workerRepo{db: db}
}

func (r *workerRepo) List(ctx context.Context, limit, offset int) ([]model.Worker, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, auth_subject, first_name, last_name, email, phone, created_at, updated_at
		FROM workers ORDER BY last_name, first_name LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list workers: %w", err)
	}
	defer rows.Close()

	var workers []model.Worker
	for rows.Next() {
		var w model.Worker
		if err := rows.Scan(&w.ID, &w.AuthSubject, &w.FirstName, &w.LastName, &w.Email, &w.Phone, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan worker: %w", err)
		}
		workers = append(workers, w)
	}
	return workers, rows.Err()
}

func (r *workerRepo) GetByID(ctx context.Context, id string) (*model.Worker, error) {
	var w model.Worker
	err := r.db.QueryRowContext(ctx,
		`SELECT id, auth_subject, first_name, last_name, email, phone, created_at, updated_at
		FROM workers WHERE id = $1`, id).
		Scan(&w.ID, &w.AuthSubject, &w.FirstName, &w.LastName, &w.Email, &w.Phone, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	return &w, nil
}

func (r *workerRepo) GetByAuthSubject(ctx context.Context, authSubject string) (*model.Worker, error) {
	var w model.Worker
	err := r.db.QueryRowContext(ctx,
		`SELECT id, auth_subject, first_name, last_name, email, phone, created_at, updated_at
		FROM workers WHERE auth_subject = $1`, authSubject).
		Scan(&w.ID, &w.AuthSubject, &w.FirstName, &w.LastName, &w.Email, &w.Phone, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get worker by auth subject: %w", err)
	}
	return &w, nil
}

func (r *workerRepo) Create(ctx context.Context, worker *model.Worker) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO workers (auth_subject, first_name, last_name, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		worker.AuthSubject, worker.FirstName, worker.LastName, worker.Email, worker.Phone).
		Scan(&worker.ID, &worker.CreatedAt, &worker.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}
	return nil
}

func (r *workerRepo) Update(ctx context.Context, worker *model.Worker) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE workers SET first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = NOW()
		WHERE id = $5`,
		worker.FirstName, worker.LastName, worker.Email, worker.Phone, worker.ID)
	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	return nil
}

func (r *workerRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM workers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete worker: %w", err)
	}
	return nil
}
