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

// ShiftReportHandler handles HTTP requests for shift reports and templates.
type ShiftReportHandler struct {
	service *service.ShiftReportService
}

// NewShiftReportHandler creates a new ShiftReportHandler.
func NewShiftReportHandler(s *service.ShiftReportService) *ShiftReportHandler {
	return &ShiftReportHandler{service: s}
}

// Routes returns the shift report routes.
func (h *ShiftReportHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/templates", func(r chi.Router) {
		// Read-only: accessible to all authenticated users
		r.Get("/", h.ListTemplates)
		r.Get("/{id}", h.GetTemplateByID)

		// Write operations: require company_admin or site_admin role
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("company_admin", "site_admin"))
			r.Post("/", h.CreateTemplate)
			r.Put("/{id}", h.UpdateTemplate)
			r.Delete("/{id}", h.DeleteTemplate)
		})
	})

	// Reports: read accessible to all, create is a worker action
	r.Get("/", h.ListReports)
	r.Get("/{id}", h.GetReportByID)
	r.Post("/", h.CreateReport)

	return r
}

func (h *ShiftReportHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("company_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	templates, err := h.service.ListTemplates(r.Context(), companyID, page, perPage)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if templates == nil {
		templates = []model.ShiftReportTemplate{}
	}
	JSON(w, http.StatusOK, templates)
}

func (h *ShiftReportHandler) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	template, err := h.service.GetTemplateByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, template)
}

func (h *ShiftReportHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var template model.ShiftReportTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.CreateTemplate(r.Context(), &template); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, template)
}

func (h *ShiftReportHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var template model.ShiftReportTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	template.ID = id

	if err := h.service.UpdateTemplate(r.Context(), &template); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusOK, template)
}

func (h *ShiftReportHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteTemplate(r.Context(), id); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShiftReportHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	shiftID := r.URL.Query().Get("shift_id")
	workerID := r.URL.Query().Get("worker_id")

	var reports []model.ShiftReport
	var err error

	if shiftID != "" {
		reports, err = h.service.ListReportsByShift(r.Context(), shiftID, page, perPage)
	} else if workerID != "" {
		reports, err = h.service.ListReportsByWorker(r.Context(), workerID, page, perPage)
	} else {
		Error(w, http.StatusBadRequest, "shift_id or worker_id query parameter is required")
		return
	}

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if reports == nil {
		reports = []model.ShiftReport{}
	}
	JSON(w, http.StatusOK, reports)
}

func (h *ShiftReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	report, err := h.service.GetReportByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, report)
}

func (h *ShiftReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var report model.ShiftReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.CreateReport(r.Context(), &report); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, report)
}
