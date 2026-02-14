package service

import (
	"context"
	"fmt"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
)

// LocationService handles business logic for location check-ins.
type LocationService struct {
	repo repository.LocationCheckInRepository
}

// NewLocationService creates a new LocationService.
func NewLocationService(repo repository.LocationCheckInRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) Create(ctx context.Context, checkIn *model.LocationCheckIn) error {
	if checkIn.WorkerID == "" {
		return fmt.Errorf("worker_id is required")
	}
	if checkIn.Latitude == 0 && checkIn.Longitude == 0 {
		return fmt.Errorf("latitude and longitude are required")
	}
	return s.repo.Create(ctx, checkIn)
}

func (s *LocationService) ListByWorker(ctx context.Context, workerID string, page, perPage int) ([]model.LocationCheckIn, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	return s.repo.ListByWorker(ctx, workerID, perPage, offset)
}

func (s *LocationService) ListByShift(ctx context.Context, shiftID string) ([]model.LocationCheckIn, error) {
	return s.repo.ListByShift(ctx, shiftID)
}
