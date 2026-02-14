package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// LocationHandler handles HTTP requests for location check-ins.
type LocationHandler struct {
	service *service.LocationService
}

// NewLocationHandler creates a new LocationHandler.
func NewLocationHandler(s *service.LocationService) *LocationHandler {
	return &LocationHandler{service: s}
}

// Routes returns the location check-in routes.
func (h *LocationHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	return r
}

func (h *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var checkIn model.LocationCheckIn
	if err := json.NewDecoder(r.Body).Decode(&checkIn); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), &checkIn); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, checkIn)
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	workerID := r.URL.Query().Get("worker_id")
	shiftID := r.URL.Query().Get("shift_id")

	var checkIns []model.LocationCheckIn
	var err error

	if workerID != "" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		checkIns, err = h.service.ListByWorker(r.Context(), workerID, page, perPage)
	} else if shiftID != "" {
		checkIns, err = h.service.ListByShift(r.Context(), shiftID)
	} else {
		Error(w, http.StatusBadRequest, "worker_id or shift_id query parameter is required")
		return
	}

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if checkIns == nil {
		checkIns = []model.LocationCheckIn{}
	}
	JSON(w, http.StatusOK, checkIns)
}
