package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// ShiftHandler handles HTTP requests for shifts.
type ShiftHandler struct {
	service *service.ShiftService
}

// NewShiftHandler creates a new ShiftHandler.
func NewShiftHandler(s *service.ShiftService) *ShiftHandler {
	return &ShiftHandler{service: s}
}

// Routes returns the shift routes.
func (h *ShiftHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Patch("/{id}/status", h.UpdateStatus)
	r.Delete("/{id}", h.Delete)
	r.Get("/{id}/assignments", h.ListAssignments)
	r.Post("/{id}/assignments", h.CreateAssignment)
	r.Patch("/{id}/assignments/{assignmentId}/accept", h.AcceptAssignment)
	r.Patch("/{id}/assignments/{assignmentId}/decline", h.DeclineAssignment)
	return r
}

func (h *ShiftHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	worksiteID := r.URL.Query().Get("worksite_id")
	status := r.URL.Query().Get("status")

	var shifts []model.Shift
	var err error

	if worksiteID != "" {
		shifts, err = h.service.ListByWorksite(r.Context(), worksiteID, page, perPage)
	} else if status != "" {
		shifts, err = h.service.ListByStatus(r.Context(), model.ShiftStatus(status), page, perPage)
	} else {
		shifts, err = h.service.List(r.Context(), page, perPage)
	}

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if shifts == nil {
		shifts = []model.Shift{}
	}
	JSON(w, http.StatusOK, shifts)
}

func (h *ShiftHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	shift, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, shift)
}

func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	var shift model.Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), &shift); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, shift)
}

func (h *ShiftHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var shift model.Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	shift.ID = id

	if err := h.service.Update(r.Context(), &shift); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusOK, shift)
}

func (h *ShiftHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		Status model.ShiftStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, body.Status); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShiftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShiftHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	assignments, err := h.service.ListAssignmentsByShift(r.Context(), id)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if assignments == nil {
		assignments = []model.ShiftAssignment{}
	}
	JSON(w, http.StatusOK, assignments)
}

func (h *ShiftHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var assignment model.ShiftAssignment
	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	assignment.ShiftID = id

	if err := h.service.CreateAssignment(r.Context(), &assignment); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, assignment)
}

func (h *ShiftHandler) AcceptAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := chi.URLParam(r, "assignmentId")

	if err := h.service.AcceptAssignment(r.Context(), assignmentID); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShiftHandler) DeclineAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := chi.URLParam(r, "assignmentId")

	if err := h.service.DeclineAssignment(r.Context(), assignmentID); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
