package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

// ShiftRepository defines the interface for shift data access.
type ShiftRepository interface {
	List(ctx context.Context, limit, offset int) ([]model.Shift, error)
	ListByWorksite(ctx context.Context, worksiteID string, limit, offset int) ([]model.Shift, error)
	ListByStatus(ctx context.Context, status model.ShiftStatus, limit, offset int) ([]model.Shift, error)
	GetByID(ctx context.Context, id string) (*model.Shift, error)
	Create(ctx context.Context, shift *model.Shift) error
	Update(ctx context.Context, shift *model.Shift) error
	UpdateStatus(ctx context.Context, id string, status model.ShiftStatus) error
	Delete(ctx context.Context, id string) error
}

type shiftRepo struct {
	db *sql.DB
}

// NewShiftRepository creates a new ShiftRepository.
func NewShiftRepository(db *sql.DB) ShiftRepository {
	return &shiftRepo{db: db}
}

func (r *shiftRepo) List(ctx context.Context, limit, offset int) ([]model.Shift, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worksite_id, created_by, title, description, start_time, end_time, status, created_at, updated_at
		 FROM shifts ORDER BY start_time DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shifts: %w", err)
	}
	defer rows.Close()

	var shifts []model.Shift
	for rows.Next() {
		var s model.Shift
		if err := rows.Scan(&s.ID, &s.WorksiteID, &s.CreatedBy, &s.Title, &s.Description, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift: %w", err)
		}
		shifts = append(shifts, s)
	}
	return shifts, rows.Err()
}

func (r *shiftRepo) ListByWorksite(ctx context.Context, worksiteID string, limit, offset int) ([]model.Shift, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worksite_id, created_by, title, description, start_time, end_time, status, created_at, updated_at
		 FROM shifts WHERE worksite_id = $1 ORDER BY start_time DESC LIMIT $2 OFFSET $3`, worksiteID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shifts by worksite: %w", err)
	}
	defer rows.Close()

	var shifts []model.Shift
	for rows.Next() {
		var s model.Shift
		if err := rows.Scan(&s.ID, &s.WorksiteID, &s.CreatedBy, &s.Title, &s.Description, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift: %w", err)
		}
		shifts = append(shifts, s)
	}
	return shifts, rows.Err()
}

func (r *shiftRepo) ListByStatus(ctx context.Context, status model.ShiftStatus, limit, offset int) ([]model.Shift, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worksite_id, created_by, title, description, start_time, end_time, status, created_at, updated_at
		 FROM shifts WHERE status = $1 ORDER BY start_time DESC LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shifts by status: %w", err)
	}
	defer rows.Close()

	var shifts []model.Shift
	for rows.Next() {
		var s model.Shift
		if err := rows.Scan(&s.ID, &s.WorksiteID, &s.CreatedBy, &s.Title, &s.Description, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shift: %w", err)
		}
		shifts = append(shifts, s)
	}
	return shifts, rows.Err()
}

func (r *shiftRepo) GetByID(ctx context.Context, id string) (*model.Shift, error) {
	var s model.Shift
	err := r.db.QueryRowContext(ctx,
		`SELECT id, worksite_id, created_by, title, description, start_time, end_time, status, created_at, updated_at
		 FROM shifts WHERE id = $1`, id).
		Scan(&s.ID, &s.WorksiteID, &s.CreatedBy, &s.Title, &s.Description, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shift: %w", err)
	}
	return &s, nil
}

func (r *shiftRepo) Create(ctx context.Context, shift *model.Shift) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO shifts (worksite_id, created_by, title, description, start_time, end_time, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		shift.WorksiteID, shift.CreatedBy, shift.Title, shift.Description, shift.StartTime, shift.EndTime, shift.Status).
		Scan(&shift.ID, &shift.CreatedAt, &shift.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create shift: %w", err)
	}
	return nil
}

func (r *shiftRepo) Update(ctx context.Context, shift *model.Shift) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE shifts SET worksite_id = $1, title = $2, description = $3, start_time = $4, end_time = $5, status = $6, updated_at = NOW()
		 WHERE id = $7`,
		shift.WorksiteID, shift.Title, shift.Description, shift.StartTime, shift.EndTime, shift.Status, shift.ID)
	if err != nil {
		return fmt.Errorf("failed to update shift: %w", err)
	}
	return nil
}

func (r *shiftRepo) UpdateStatus(ctx context.Context, id string, status model.ShiftStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE shifts SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update shift status: %w", err)
	}
	return nil
}

func (r *shiftRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shifts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete shift: %w", err)
	}
	return nil
}
