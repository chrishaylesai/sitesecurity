package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// mockShiftRepo is a test double for repository.ShiftRepository.
type mockShiftRepo struct {
	shifts []model.Shift
	err    error
}

func (m *mockShiftRepo) List(ctx context.Context, limit, offset int) ([]model.Shift, error) {
	if m.err != nil {
		return nil, m.err
	}
	end := offset + limit
	if end > len(m.shifts) {
		end = len(m.shifts)
	}
	if offset >= len(m.shifts) {
		return nil, nil
	}
	return m.shifts[offset:end], nil
}

func (m *mockShiftRepo) ListByWorksite(ctx context.Context, worksiteID string, limit, offset int) ([]model.Shift, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.Shift
	for _, s := range m.shifts {
		if s.WorksiteID == worksiteID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockShiftRepo) ListByStatus(ctx context.Context, status model.ShiftStatus, limit, offset int) ([]model.Shift, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.Shift
	for _, s := range m.shifts {
		if s.Status == status {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockShiftRepo) GetByID(ctx context.Context, id string) (*model.Shift, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, s := range m.shifts {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, nil
}

func (m *mockShiftRepo) Create(ctx context.Context, shift *model.Shift) error {
	if m.err != nil {
		return m.err
	}
	shift.ID = "new-shift-id"
	m.shifts = append(m.shifts, *shift)
	return nil
}

func (m *mockShiftRepo) Update(ctx context.Context, shift *model.Shift) error {
	return m.err
}

func (m *mockShiftRepo) UpdateStatus(ctx context.Context, id string, status model.ShiftStatus) error {
	return m.err
}

func (m *mockShiftRepo) Delete(ctx context.Context, id string) error {
	return m.err
}

// mockShiftAssignmentRepo is a test double for repository.ShiftAssignmentRepository.
type mockShiftAssignmentRepo struct {
	assignments []model.ShiftAssignment
	err         error
}

func (m *mockShiftAssignmentRepo) ListByShift(ctx context.Context, shiftID string) ([]model.ShiftAssignment, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.ShiftAssignment
	for _, a := range m.assignments {
		if a.ShiftID == shiftID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockShiftAssignmentRepo) ListByWorker(ctx context.Context, workerID string) ([]model.ShiftAssignment, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []model.ShiftAssignment
	for _, a := range m.assignments {
		if a.WorkerID == workerID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockShiftAssignmentRepo) GetByID(ctx context.Context, id string) (*model.ShiftAssignment, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, a := range m.assignments {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, nil
}

func (m *mockShiftAssignmentRepo) Get(ctx context.Context, shiftID, workerID string) (*model.ShiftAssignment, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, a := range m.assignments {
		if a.ShiftID == shiftID && a.WorkerID == workerID {
			return &a, nil
		}
	}
	return nil, nil
}

func (m *mockShiftAssignmentRepo) Create(ctx context.Context, assignment *model.ShiftAssignment) error {
	if m.err != nil {
		return m.err
	}
	assignment.ID = "new-assignment-id"
	m.assignments = append(m.assignments, *assignment)
	return nil
}

func (m *mockShiftAssignmentRepo) UpdateStatus(ctx context.Context, id string, status model.AssignmentStatus) error {
	return m.err
}

func (m *mockShiftAssignmentRepo) Delete(ctx context.Context, id string) error {
	return m.err
}

func TestShiftService_Create_Valid(t *testing.T) {
	shiftRepo := &mockShiftRepo{}
	assignmentRepo := &mockShiftAssignmentRepo{}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	now := time.Now()
	shift := &model.Shift{
		Title:      "Night Watch",
		WorksiteID: "ws-1",
		StartTime:  now,
		EndTime:    now.Add(8 * time.Hour),
	}
	err := svc.Create(context.Background(), shift)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shift.ID != "new-shift-id" {
		t.Errorf("expected ID to be set, got '%s'", shift.ID)
	}
	if shift.Status != model.ShiftOpen {
		t.Errorf("expected status 'open', got '%s'", shift.Status)
	}
}

func TestShiftService_Create_MissingTitle(t *testing.T) {
	shiftRepo := &mockShiftRepo{}
	assignmentRepo := &mockShiftAssignmentRepo{}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	now := time.Now()
	shift := &model.Shift{
		Title:      "",
		WorksiteID: "ws-1",
		StartTime:  now,
		EndTime:    now.Add(8 * time.Hour),
	}
	err := svc.Create(context.Background(), shift)
	if err == nil {
		t.Error("expected error for missing title")
	}
}

func TestShiftService_Create_InvalidTimeRange(t *testing.T) {
	shiftRepo := &mockShiftRepo{}
	assignmentRepo := &mockShiftAssignmentRepo{}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	now := time.Now()
	shift := &model.Shift{
		Title:      "Night Watch",
		WorksiteID: "ws-1",
		StartTime:  now.Add(8 * time.Hour),
		EndTime:    now,
	}
	err := svc.Create(context.Background(), shift)
	if err == nil {
		t.Error("expected error for invalid time range")
	}
}

func TestShiftService_GetByID_NotFound(t *testing.T) {
	shiftRepo := &mockShiftRepo{shifts: []model.Shift{}}
	assignmentRepo := &mockShiftAssignmentRepo{}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	_, err := svc.GetByID(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing shift")
	}
}

func TestShiftService_AcceptAssignment_Valid(t *testing.T) {
	shiftRepo := &mockShiftRepo{}
	assignmentRepo := &mockShiftAssignmentRepo{
		assignments: []model.ShiftAssignment{
			{ID: "a-1", ShiftID: "s-1", WorkerID: "w-1", Status: model.AssignmentOffered},
		},
	}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	err := svc.AcceptAssignment(context.Background(), "a-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShiftService_AcceptAssignment_NotFound(t *testing.T) {
	shiftRepo := &mockShiftRepo{}
	assignmentRepo := &mockShiftAssignmentRepo{assignments: []model.ShiftAssignment{}}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	err := svc.AcceptAssignment(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing assignment")
	}
}

func TestShiftService_DeclineAssignment_Valid(t *testing.T) {
	shiftRepo := &mockShiftRepo{}
	assignmentRepo := &mockShiftAssignmentRepo{
		assignments: []model.ShiftAssignment{
			{ID: "a-1", ShiftID: "s-1", WorkerID: "w-1", Status: model.AssignmentOffered},
		},
	}
	svc := service.NewShiftService(shiftRepo, assignmentRepo)

	err := svc.DeclineAssignment(context.Background(), "a-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
