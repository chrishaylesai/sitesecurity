package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

// ShiftService handles business logic for shifts and shift assignments.
type ShiftService struct {
	shiftRepo      repository.ShiftRepository
	assignmentRepo repository.ShiftAssignmentRepository
}

// NewShiftService creates a new ShiftService.
func NewShiftService(shiftRepo repository.ShiftRepository, assignmentRepo repository.ShiftAssignmentRepository) *ShiftService {
	return &ShiftService{shiftRepo: shiftRepo, assignmentRepo: assignmentRepo}
}

func (s *ShiftService) List(ctx context.Context, page, perPage int) ([]model.Shift, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.shiftRepo.List(ctx, perPage, offset)
}

func (s *ShiftService) ListByWorksite(ctx context.Context, worksiteID string, page, perPage int) ([]model.Shift, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.shiftRepo.ListByWorksite(ctx, worksiteID, perPage, offset)
}

func (s *ShiftService) ListByStatus(ctx context.Context, status model.ShiftStatus, page, perPage int) ([]model.Shift, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.shiftRepo.ListByStatus(ctx, status, perPage, offset)
}

func (s *ShiftService) GetByID(ctx context.Context, id string) (*model.Shift, error) {
	shift, err := s.shiftRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if shift == nil {
		return nil, fmt.Errorf("shift not found")
	}
	return shift, nil
}

func (s *ShiftService) Create(ctx context.Context, shift *model.Shift) error {
	if shift.Title == "" {
		return fmt.Errorf("shift title is required")
	}
	if shift.WorksiteID == "" {
		return fmt.Errorf("worksite_id is required")
	}
	if !shift.StartTime.Before(shift.EndTime) {
		return fmt.Errorf("start_time must be before end_time")
	}
	if shift.Status == "" {
		shift.Status = model.ShiftOpen
	}
	return s.shiftRepo.Create(ctx, shift)
}

func (s *ShiftService) Update(ctx context.Context, shift *model.Shift) error {
	if shift.Title == "" {
		return fmt.Errorf("shift title is required")
	}
	if shift.WorksiteID == "" {
		return fmt.Errorf("worksite_id is required")
	}
	if !shift.StartTime.Before(shift.EndTime) {
		return fmt.Errorf("start_time must be before end_time")
	}
	existing, err := s.shiftRepo.GetByID(ctx, shift.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("shift not found")
	}
	return s.shiftRepo.Update(ctx, shift)
}

func (s *ShiftService) UpdateStatus(ctx context.Context, id string, status model.ShiftStatus) error {
	existing, err := s.shiftRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("shift not found")
	}
	if !isValidShiftTransition(existing.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", existing.Status, status)
	}
	return s.shiftRepo.UpdateStatus(ctx, id, status)
}

func (s *ShiftService) Delete(ctx context.Context, id string) error {
	existing, err := s.shiftRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("shift not found")
	}
	return s.shiftRepo.Delete(ctx, id)
}

// isValidShiftTransition checks whether a shift status transition is allowed.
func isValidShiftTransition(from, to model.ShiftStatus) bool {
	if to == model.ShiftCancelled {
		return true
	}
	switch from {
	case model.ShiftOpen:
		return to == model.ShiftAssigned
	case model.ShiftAssigned:
		return to == model.ShiftInProgress
	case model.ShiftInProgress:
		return to == model.ShiftCompleted
	default:
		return false
	}
}

// CreateAssignment creates a new shift assignment (offers a shift to a worker).
func (s *ShiftService) CreateAssignment(ctx context.Context, assignment *model.ShiftAssignment) error {
	shift, err := s.shiftRepo.GetByID(ctx, assignment.ShiftID)
	if err != nil {
		return err
	}
	if shift == nil {
		return fmt.Errorf("shift not found")
	}
	if assignment.Status == "" {
		assignment.Status = model.AssignmentOffered
	}
	return s.assignmentRepo.Create(ctx, assignment)
}

// AcceptAssignment marks a shift assignment as accepted.
func (s *ShiftService) AcceptAssignment(ctx context.Context, id string) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return fmt.Errorf("assignment not found")
	}
	if assignment.Status != model.AssignmentOffered {
		return fmt.Errorf("assignment cannot be accepted from status %s", assignment.Status)
	}
	return s.assignmentRepo.UpdateStatus(ctx, id, model.AssignmentAccepted)
}

// DeclineAssignment marks a shift assignment as declined.
func (s *ShiftService) DeclineAssignment(ctx context.Context, id string) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return fmt.Errorf("assignment not found")
	}
	if assignment.Status != model.AssignmentOffered {
		return fmt.Errorf("assignment cannot be declined from status %s", assignment.Status)
	}
	return s.assignmentRepo.UpdateStatus(ctx, id, model.AssignmentDeclined)
}

// CompleteAssignment marks a shift assignment as completed.
func (s *ShiftService) CompleteAssignment(ctx context.Context, id string) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return fmt.Errorf("assignment not found")
	}
	if assignment.Status != model.AssignmentAccepted {
		return fmt.Errorf("assignment cannot be completed from status %s", assignment.Status)
	}
	return s.assignmentRepo.UpdateStatus(ctx, id, model.AssignmentCompleted)
}

// ListAssignmentsByShift returns all assignments for a given shift.
func (s *ShiftService) ListAssignmentsByShift(ctx context.Context, shiftID string) ([]model.ShiftAssignment, error) {
	return s.assignmentRepo.ListByShift(ctx, shiftID)
}
