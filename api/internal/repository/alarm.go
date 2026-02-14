package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
)

// AlarmRepository defines the interface for alarm data access.
type AlarmRepository interface {
	List(ctx context.Context, limit, offset int) ([]model.Alarm, error)
	ListByStatus(ctx context.Context, status model.AlarmStatus, limit, offset int) ([]model.Alarm, error)
	ListByWorker(ctx context.Context, workerID string) ([]model.Alarm, error)
	GetByID(ctx context.Context, id string) (*model.Alarm, error)
	Create(ctx context.Context, alarm *model.Alarm) error
	UpdateStatus(ctx context.Context, id string, status model.AlarmStatus) error
}

type alarmRepo struct {
	db *sql.DB
}

// NewAlarmRepository creates a new AlarmRepository.
func NewAlarmRepository(db *sql.DB) AlarmRepository {
	return &alarmRepo{db: db}
}

func (r *alarmRepo) List(ctx context.Context, limit, offset int) ([]model.Alarm, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worker_id, shift_id, latitude, longitude, message, status, raised_at, acknowledged_at, resolved_at
		 FROM alarms ORDER BY raised_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list alarms: %w", err)
	}
	defer rows.Close()

	var alarms []model.Alarm
	for rows.Next() {
		var a model.Alarm
		if err := rows.Scan(&a.ID, &a.WorkerID, &a.ShiftID, &a.Latitude, &a.Longitude, &a.Message, &a.Status, &a.RaisedAt, &a.AcknowledgedAt, &a.ResolvedAt); err != nil {
			return nil, fmt.Errorf("failed to scan alarm: %w", err)
		}
		alarms = append(alarms, a)
	}
	return alarms, rows.Err()
}

func (r *alarmRepo) ListByStatus(ctx context.Context, status model.AlarmStatus, limit, offset int) ([]model.Alarm, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worker_id, shift_id, latitude, longitude, message, status, raised_at, acknowledged_at, resolved_at
		 FROM alarms WHERE status = $1 ORDER BY raised_at DESC LIMIT $2 OFFSET $3`,
		status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list alarms by status: %w", err)
	}
	defer rows.Close()

	var alarms []model.Alarm
	for rows.Next() {
		var a model.Alarm
		if err := rows.Scan(&a.ID, &a.WorkerID, &a.ShiftID, &a.Latitude, &a.Longitude, &a.Message, &a.Status, &a.RaisedAt, &a.AcknowledgedAt, &a.ResolvedAt); err != nil {
			return nil, fmt.Errorf("failed to scan alarm: %w", err)
		}
		alarms = append(alarms, a)
	}
	return alarms, rows.Err()
}

func (r *alarmRepo) ListByWorker(ctx context.Context, workerID string) ([]model.Alarm, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, worker_id, shift_id, latitude, longitude, message, status, raised_at, acknowledged_at, resolved_at
		 FROM alarms WHERE worker_id = $1 ORDER BY raised_at DESC`,
		workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list alarms by worker: %w", err)
	}
	defer rows.Close()

	var alarms []model.Alarm
	for rows.Next() {
		var a model.Alarm
		if err := rows.Scan(&a.ID, &a.WorkerID, &a.ShiftID, &a.Latitude, &a.Longitude, &a.Message, &a.Status, &a.RaisedAt, &a.AcknowledgedAt, &a.ResolvedAt); err != nil {
			return nil, fmt.Errorf("failed to scan alarm: %w", err)
		}
		alarms = append(alarms, a)
	}
	return alarms, rows.Err()
}

func (r *alarmRepo) GetByID(ctx context.Context, id string) (*model.Alarm, error) {
	var a model.Alarm
	err := r.db.QueryRowContext(ctx,
		`SELECT id, worker_id, shift_id, latitude, longitude, message, status, raised_at, acknowledged_at, resolved_at
		 FROM alarms WHERE id = $1`, id).
		Scan(&a.ID, &a.WorkerID, &a.ShiftID, &a.Latitude, &a.Longitude, &a.Message, &a.Status, &a.RaisedAt, &a.AcknowledgedAt, &a.ResolvedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alarm: %w", err)
	}
	return &a, nil
}

func (r *alarmRepo) Create(ctx context.Context, alarm *model.Alarm) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO alarms (worker_id, shift_id, latitude, longitude, message, status)
		 VALUES ($1, $2, $3, $4, $5, 'raised')
		 RETURNING id, raised_at`,
		alarm.WorkerID, alarm.ShiftID, alarm.Latitude, alarm.Longitude, alarm.Message).
		Scan(&alarm.ID, &alarm.RaisedAt)
	if err != nil {
		return fmt.Errorf("failed to create alarm: %w", err)
	}
	alarm.Status = model.AlarmRaised
	return nil
}

func (r *alarmRepo) UpdateStatus(ctx context.Context, id string, status model.AlarmStatus) error {
	var query string
	switch status {
	case model.AlarmAcknowledged:
		query = `UPDATE alarms SET status = $1, acknowledged_at = $2 WHERE id = $3`
	case model.AlarmResolved:
		query = `UPDATE alarms SET status = $1, resolved_at = $2 WHERE id = $3`
	default:
		return fmt.Errorf("unsupported status transition: %s", status)
	}

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update alarm status: %w", err)
	}
	return nil
}
