package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/middleware"
	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// AlarmHandler handles HTTP requests for alarms.
type AlarmHandler struct {
	service *service.AlarmService
}

// NewAlarmHandler creates a new AlarmHandler.
func NewAlarmHandler(s *service.AlarmService) *AlarmHandler {
	return &AlarmHandler{service: s}
}

// Routes returns the alarm routes.
func (h *AlarmHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Read-only: accessible to all authenticated users
	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)

	// Worker-specific actions: accessible to all authenticated users
	r.Post("/", h.Raise)

	// Admin actions: require company_admin or site_admin role
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("company_admin", "site_admin"))
		r.Patch("/{id}/acknowledge", h.Acknowledge)
		r.Patch("/{id}/resolve", h.Resolve)
	})

	return r
}

func (h *AlarmHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	status := r.URL.Query().Get("status")
	workerID := r.URL.Query().Get("worker_id")

	var alarms []model.Alarm
	var err error

	if workerID != "" {
		alarms, err = h.service.ListByWorker(r.Context(), workerID)
	} else if status != "" {
		alarms, err = h.service.ListByStatus(r.Context(), model.AlarmStatus(status), page, perPage)
	} else {
		alarms, err = h.service.List(r.Context(), page, perPage)
	}

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if alarms == nil {
		alarms = []model.Alarm{}
	}
	JSON(w, http.StatusOK, alarms)
}

func (h *AlarmHandler) Raise(w http.ResponseWriter, r *http.Request) {
	var alarm model.Alarm
	if err := json.NewDecoder(r.Body).Decode(&alarm); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Raise(r.Context(), &alarm); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, alarm)
}

func (h *AlarmHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	alarm, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, alarm)
}

func (h *AlarmHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Acknowledge(r.Context(), id); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
}

func (h *AlarmHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Resolve(r.Context(), id); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"status": "resolved"})
}
