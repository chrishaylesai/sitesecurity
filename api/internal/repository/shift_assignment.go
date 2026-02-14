package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

// ShiftAssignmentRepository defines the interface for shift assignment data access.
type ShiftAssignmentRepository interface {
	ListByShift(ctx context.Context, shiftID string) ([]model.ShiftAssignment, error)
	ListByWorker(ctx context.Context, workerID string) ([]model.ShiftAssignment, error)
	GetByID(ctx context.Context, id string) (*model.ShiftAssignment, error)
	Get(ctx context.Context, shiftID, workerID string) (*model.ShiftAssignment, error)
	Create(ctx context.Context, assignment *model.ShiftAssignment) error
	UpdateStatus(ctx context.Context, id string, status model.AssignmentStatus) error
	Delete(ctx context.Context, id string) error
}

type shiftAssignmentRepo struct {
	db *sql.DB
}

// NewShiftAssignmentRepository creates a new ShiftAssignmentRepository.
func NewShiftAssignmentRepository(db *sql.DB) ShiftAssignmentRepository {
	return &shiftAssignmentRepo{db: db}
}

func (r *shiftAssignmentRepo) ListByShift(ctx context.Context, shiftID string) ([]model.ShiftAssignment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shift_id, worker_id, status, assigned_at, responded_at
		 FROM shift_assignments WHERE shift_id = $1 ORDER BY assigned_at DESC`, shiftID)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignments by shift: %w", err)
	}
	defer rows.Close()

	var assignments []model.ShiftAssignment
	for rows.Next() {
		var a model.ShiftAssignment
		if err := rows.Scan(&a.ID, &a.ShiftID, &a.WorkerID, &a.Status, &a.AssignedAt, &a.RespondedAt); err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func (r *shiftAssignmentRepo) ListByWorker(ctx context.Context, workerID string) ([]model.ShiftAssignment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shift_id, worker_id, status, assigned_at, responded_at
		 FROM shift_assignments WHERE worker_id = $1 ORDER BY assigned_at DESC`, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignments by worker: %w", err)
	}
	defer rows.Close()

	var assignments []model.ShiftAssignment
	for rows.Next() {
		var a model.ShiftAssignment
		if err := rows.Scan(&a.ID, &a.ShiftID, &a.WorkerID, &a.Status, &a.AssignedAt, &a.RespondedAt); err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func (r *shiftAssignmentRepo) GetByID(ctx context.Context, id string) (*model.ShiftAssignment, error) {
	var a model.ShiftAssignment
	err := r.db.QueryRowContext(ctx,
		`SELECT id, shift_id, worker_id, status, assigned_at, responded_at
		 FROM shift_assignments WHERE id = $1`, id).
		Scan(&a.ID, &a.ShiftID, &a.WorkerID, &a.Status, &a.AssignedAt, &a.RespondedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}
	return &a, nil
}

func (r *shiftAssignmentRepo) Get(ctx context.Context, shiftID, workerID string) (*model.ShiftAssignment, error) {
	var a model.ShiftAssignment
	err := r.db.QueryRowContext(ctx,
		`SELECT id, shift_id, worker_id, status, assigned_at, responded_at
		 FROM shift_assignments WHERE shift_id = $1 AND worker_id = $2`, shiftID, workerID).
		Scan(&a.ID, &a.ShiftID, &a.WorkerID, &a.Status, &a.AssignedAt, &a.RespondedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}
	return &a, nil
}

func (r *shiftAssignmentRepo) Create(ctx context.Context, assignment *model.ShiftAssignment) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO shift_assignments (shift_id, worker_id, status)
		 VALUES ($1, $2, $3)
		 RETURNING id, assigned_at`,
		assignment.ShiftID, assignment.WorkerID, assignment.Status).
		Scan(&assignment.ID, &assignment.AssignedAt)
	if err != nil {
		return fmt.Errorf("failed to create assignment: %w", err)
	}
	return nil
}

func (r *shiftAssignmentRepo) UpdateStatus(ctx context.Context, id string, status model.AssignmentStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE shift_assignments SET status = $1, responded_at = NOW() WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update assignment status: %w", err)
	}
	return nil
}

func (r *shiftAssignmentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shift_assignments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}
	return nil
}
