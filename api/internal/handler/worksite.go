package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/model"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

type WorksiteHandler struct {
	service *service.WorksiteService
}

func NewWorksiteHandler(s *service.WorksiteService) *WorksiteHandler {
	return &WorksiteHandler{service: s}
}

func (h *WorksiteHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *WorksiteHandler) List(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("company_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	worksites, err := h.service.List(r.Context(), companyID, page, perPage)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if worksites == nil {
		worksites = []model.Worksite{}
	}
	JSON(w, http.StatusOK, worksites)
}

func (h *WorksiteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	worksite, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, worksite)
}

func (h *WorksiteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var worksite model.Worksite
	if err := json.NewDecoder(r.Body).Decode(&worksite); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.Create(r.Context(), &worksite); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusCreated, worksite)
}

func (h *WorksiteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var worksite model.Worksite
	if err := json.NewDecoder(r.Body).Decode(&worksite); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	worksite.ID = id
	if err := h.service.Update(r.Context(), &worksite); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusOK, worksite)
}

func (h *WorksiteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
