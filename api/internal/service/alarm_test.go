package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// mockAlarmRepo is a test double for repository.AlarmRepository.
type mockAlarmRepo struct {
	alarms []model.Alarm
	err    error
}

func (m *mockAlarmRepo) List(ctx context.Context, limit, offset int) ([]model.Alarm, error) {
	if m.err != nil {
		return nil, m.err
	}
	end := offset + limit
	if end > len(m.alarms) {
		end = len(m.alarms)
	}
	if offset >= len(m.alarms) {
		return nil, nil
	}
	return m.alarms[offset:end], nil
}

func (m *mockAlarmRepo) ListByStatus(ctx context.Context, status model.AlarmStatus, limit, offset int) ([]model.Alarm, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.Alarm
	for _, a := range m.alarms {
		if a.Status == status {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockAlarmRepo) ListByWorker(ctx context.Context, workerID string) ([]model.Alarm, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.Alarm
	for _, a := range m.alarms {
		if a.WorkerID == workerID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockAlarmRepo) GetByID(ctx context.Context, id string) (*model.Alarm, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, a := range m.alarms {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, nil
}

func (m *mockAlarmRepo) Create(ctx context.Context, alarm *model.Alarm) error {
	if m.err != nil {
		return m.err
	}
	alarm.ID = "new-alarm-id"
	alarm.Status = model.AlarmRaised
	alarm.RaisedAt = time.Now()
	m.alarms = append(m.alarms, *alarm)
	return nil
}

func (m *mockAlarmRepo) UpdateStatus(ctx context.Context, id string, status model.AlarmStatus) error {
	if m.err != nil {
		return m.err
	}
	for i := range m.alarms {
		if m.alarms[i].ID == id {
			m.alarms[i].Status = status
			return nil
		}
	}
	return nil
}

func TestAlarmService_Raise_Valid(t *testing.T) {
	repo := &mockAlarmRepo{}
	svc := service.NewAlarmService(repo)

	alarm := &model.Alarm{WorkerID: "worker-1"}
	err := svc.Raise(context.Background(), alarm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alarm.ID != "new-alarm-id" {
		t.Errorf("expected ID to be set, got '%s'", alarm.ID)
	}
}

func TestAlarmService_Raise_MissingWorkerID(t *testing.T) {
	repo := &mockAlarmRepo{}
	svc := service.NewAlarmService(repo)

	alarm := &model.Alarm{}
	err := svc.Raise(context.Background(), alarm)
	if err == nil {
		t.Error("expected error for missing worker_id")
	}
}

func TestAlarmService_Acknowledge_Valid(t *testing.T) {
	repo := &mockAlarmRepo{
		alarms: []model.Alarm{
			{ID: "alarm-1", WorkerID: "worker-1", Status: model.AlarmRaised},
		},
	}
	svc := service.NewAlarmService(repo)

	err := svc.Acknowledge(context.Background(), "alarm-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlarmService_Acknowledge_WrongStatus(t *testing.T) {
	repo := &mockAlarmRepo{
		alarms: []model.Alarm{
			{ID: "alarm-1", WorkerID: "worker-1", Status: model.AlarmResolved},
		},
	}
	svc := service.NewAlarmService(repo)

	err := svc.Acknowledge(context.Background(), "alarm-1")
	if err == nil {
		t.Error("expected error for wrong status")
	}
}

func TestAlarmService_Resolve_Valid(t *testing.T) {
	repo := &mockAlarmRepo{
		alarms: []model.Alarm{
			{ID: "alarm-1", WorkerID: "worker-1", Status: model.AlarmAcknowledged},
		},
	}
	svc := service.NewAlarmService(repo)

	err := svc.Resolve(context.Background(), "alarm-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlarmService_GetByID_NotFound(t *testing.T) {
	repo := &mockAlarmRepo{alarms: []model.Alarm{}}
	svc := service.NewAlarmService(repo)

	_, err := svc.GetByID(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing alarm")
	}
}
