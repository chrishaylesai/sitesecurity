package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

// LocationCheckInRepository defines the interface for location check-in data access.
type LocationCheckInRepository interface {
	ListByWorker(ctx context.Context, workerID string, limit, offset int) ([]model.LocationCheckIn, error)
	ListByShift(ctx context.Context, shiftID string) ([]model.LocationCheckIn, error)
	Create(ctx context.Context, checkIn *model.LocationCheckIn) error
}

type locationCheckInRepo struct {
	db *sql.DB
}

// NewLocationCheckInRepository creates a new LocationCheckInRepository.
func NewLocationCheckInRepository(db *sql.DB) LocationCheckInRepository {
	return &locationCheckInRepo{db: db}
}

func (r *locationCheckInRepo) ListByWorker(ctx context.Context, workerID string, limit, offset int) ([]model.LocationCheckIn, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worker_id, shift_id, latitude, longitude, recorded_at
		 FROM location_check_ins WHERE worker_id = $1 ORDER BY recorded_at DESC LIMIT $2 OFFSET $3`,
		workerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list location check-ins by worker: %w", err)
	}
	defer rows.Close()

	var checkIns []model.LocationCheckIn
	for rows.Next() {
		var c model.LocationCheckIn
		if err := rows.Scan(&c.ID, &c.WorkerID, &c.ShiftID, &c.Latitude, &c.Longitude, &c.RecordedAt); err != nil {
			return nil, fmt.Errorf("failed to scan location check-in: %w", err)
		}
		checkIns = append(checkIns, c)
	}
	return checkIns, rows.Err()
}

func (r *locationCheckInRepo) ListByShift(ctx context.Context, shiftID string) ([]model.LocationCheckIn, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worker_id, shift_id, latitude, longitude, recorded_at
		 FROM location_check_ins WHERE shift_id = $1 ORDER BY recorded_at DESC`,
		shiftID)
	if err != nil {
		return nil, fmt.Errorf("failed to list location check-ins by shift: %w", err)
	}
	defer rows.Close()

	var checkIns []model.LocationCheckIn
	for rows.Next() {
		var c model.LocationCheckIn
		if err := rows.Scan(&c.ID, &c.WorkerID, &c.ShiftID, &c.Latitude, &c.Longitude, &c.RecordedAt); err != nil {
			return nil, fmt.Errorf("failed to scan location check-in: %w", err)
		}
		checkIns = append(checkIns, c)
	}
	return checkIns, rows.Err()
}

func (r *locationCheckInRepo) Create(ctx context.Context, checkIn *model.LocationCheckIn) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO location_check_ins (worker_id, shift_id, latitude, longitude)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, recorded_at`,
		checkIn.WorkerID, checkIn.ShiftID, checkIn.Latitude, checkIn.Longitude).
		Scan(&checkIn.ID, &checkIn.RecordedAt)
	if err != nil {
		return fmt.Errorf("failed to create location check-in: %w", err)
	}
	return nil
}
