package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

// AlarmService handles business logic for alarms.
type AlarmService struct {
	repo repository.AlarmRepository
}

// NewAlarmService creates a new AlarmService.
func NewAlarmService(repo repository.AlarmRepository) *AlarmService {
	return &AlarmService{repo: repo}
}

func (s *AlarmService) List(ctx context.Context, page, perPage int) ([]model.Alarm, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.repo.List(ctx, perPage, offset)
}

func (s *AlarmService) ListByStatus(ctx context.Context, status model.AlarmStatus, page, perPage int) ([]model.Alarm, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.repo.ListByStatus(ctx, status, perPage, offset)
}

func (s *AlarmService) ListByWorker(ctx context.Context, workerID string) ([]model.Alarm, error) {
	return s.repo.ListByWorker(ctx, workerID)
}

func (s *AlarmService) GetByID(ctx context.Context, id string) (*model.Alarm, error) {
	alarm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if alarm == nil {
		return nil, fmt.Errorf("alarm not found")
	}
	return alarm, nil
}

func (s *AlarmService) Raise(ctx context.Context, alarm *model.Alarm) error {
	if alarm.WorkerID == "" {
		return fmt.Errorf("worker_id is required")
	}
	return s.repo.Create(ctx, alarm)
}

func (s *AlarmService) Acknowledge(ctx context.Context, id string) error {
	alarm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if alarm == nil {
		return fmt.Errorf("alarm not found")
	}
	if alarm.Status != model.AlarmRaised {
		return fmt.Errorf("alarm must be in 'raised' status to acknowledge")
	}
	return s.repo.UpdateStatus(ctx, id, model.AlarmAcknowledged)
}

func (s *AlarmService) Resolve(ctx context.Context, id string) error {
	alarm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if alarm == nil {
		return fmt.Errorf("alarm not found")
	}
	if alarm.Status != model.AlarmRaised && alarm.Status != model.AlarmAcknowledged {
		return fmt.Errorf("alarm must be in 'raised' or 'acknowledged' status to resolve")
	}
	return s.repo.UpdateStatus(ctx, id, model.AlarmResolved)
}
