package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

// CompanyHandler handles HTTP requests for companies.
type CompanyHandler struct {
	service *service.CompanyService
}

// NewCompanyHandler creates a new CompanyHandler.
func NewCompanyHandler(s *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: s}
}

// Routes returns the company routes.
func (h *CompanyHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	companies, err := h.service.List(r.Context(), page, perPage)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if companies == nil {
		companies = []model.Company{}
	}
	JSON(w, http.StatusOK, companies)
}

func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	company, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, company)
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), &company); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusCreated, company)
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var company model.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	company.ID = id

	if err := h.service.Update(r.Context(), &company); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	JSON(w, http.StatusOK, company)
}

func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
