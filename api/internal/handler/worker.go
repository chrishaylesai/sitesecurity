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

type WorkerHandler struct {
	service *service.WorkerService
}

func NewWorkerHandler(s *service.WorkerService) *WorkerHandler {
	return &WorkerHandler{service: s}
}

func (h *WorkerHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Read-only: accessible to all authenticated users
	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)
	r.Get("/{id}/certificates", h.ListCertificates)
	r.Get("/{id}/certificates/{certId}", h.GetCertificate)
	r.Get("/{id}/memberships", h.ListMemberships)

	// Write operations: require company_admin role
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("company_admin"))
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)

		// Certificates
		r.Post("/{id}/certificates", h.CreateCertificate)
		r.Put("/{id}/certificates/{certId}", h.UpdateCertificate)
		r.Delete("/{id}/certificates/{certId}", h.DeleteCertificate)

		// Memberships
		r.Post("/{id}/memberships", h.AddMembership)
		r.Put("/{id}/memberships/{companyId}", h.UpdateMembershipRole)
		r.Delete("/{id}/memberships/{companyId}", h.RemoveMembership)
	})

	return r
}

func (h *WorkerHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	workers, err := h.service.List(r.Context(), page, perPage)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if workers == nil {
		workers = []model.Worker{}
	}
	JSON(w, http.StatusOK, workers)
}

func (h *WorkerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	worker, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, worker)
}

func (h *WorkerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var worker model.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.Create(r.Context(), &worker); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusCreated, worker)
}

func (h *WorkerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var worker model.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	worker.ID = id
	if err := h.service.Update(r.Context(), &worker); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusOK, worker)
}

// Certificates

func (h *WorkerHandler) ListCertificates(w http.ResponseWriter, r *http.Request) {
	workerID := chi.URLParam(r, "id")
	certs, err := h.service.ListCertificates(r.Context(), workerID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if certs == nil {
		certs = []model.Certificate{}
	}
	JSON(w, http.StatusOK, certs)
}

func (h *WorkerHandler) GetCertificate(w http.ResponseWriter, r *http.Request) {
	certID := chi.URLParam(r, "certId")
	cert, err := h.service.GetCertificate(r.Context(), certID)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, cert)
}

func (h *WorkerHandler) CreateCertificate(w http.ResponseWriter, r *http.Request) {
	workerID := chi.URLParam(r, "id")
	var cert model.Certificate
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cert.WorkerID = workerID
	if err := h.service.CreateCertificate(r.Context(), &cert); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusCreated, cert)
}

func (h *WorkerHandler) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	certID := chi.URLParam(r, "certId")
	var cert model.Certificate
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cert.ID = certID
	if err := h.service.UpdateCertificate(r.Context(), &cert); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusOK, cert)
}

func (h *WorkerHandler) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	certID := chi.URLParam(r, "certId")
	if err := h.service.DeleteCertificate(r.Context(), certID); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Memberships

func (h *WorkerHandler) ListMemberships(w http.ResponseWriter, r *http.Request) {
	workerID := chi.URLParam(r, "id")
	memberships, err := h.service.ListMemberships(r.Context(), workerID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if memberships == nil {
		memberships = []model.WorkerCompany{}
	}
	JSON(w, http.StatusOK, memberships)
}

func (h *WorkerHandler) AddMembership(w http.ResponseWriter, r *http.Request) {
	workerID := chi.URLParam(r, "id")
	var wc model.WorkerCompany
	if err := json.NewDecoder(r.Body).Decode(&wc); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	wc.WorkerID = workerID
	if err := h.service.AddMembership(r.Context(), &wc); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	JSON(w, http.StatusCreated, wc)
}

func (h *WorkerHandler) UpdateMembershipRole(w http.ResponseWriter, r *http.Request) {
	workerID := chi.URLParam(r, "id")
	companyID := chi.URLParam(r, "companyId")

	var body struct {
		Role model.WorkerRole `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.UpdateMembershipRole(r.Context(), workerID, companyID, body.Role); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkerHandler) RemoveMembership(w http.ResponseWriter, r *http.Request) {
	workerID := chi.URLParam(r, "id")
	companyID := chi.URLParam(r, "companyId")
	if err := h.service.RemoveMembership(r.Context(), workerID, companyID); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
