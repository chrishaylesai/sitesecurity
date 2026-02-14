package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/chrishaylesai/sitesecurity/api/internal/auth/keycloak"
	"github.com/chrishaylesai/sitesecurity/api/internal/config"
	"github.com/chrishaylesai/sitesecurity/api/internal/handler"
	"github.com/chrishaylesai/sitesecurity/api/internal/middleware"
	"github.com/chrishaylesai/sitesecurity/api/internal/repository"
	"github.com/chrishaylesai/sitesecurity/api/internal/service"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auth provider
	authProvider := keycloak.New(cfg.Auth)

	// Repositories
	companyRepo := repository.NewCompanyRepository(db)
	worksiteRepo := repository.NewWorksiteRepository(db)
	workerRepo := repository.NewWorkerRepository(db)
	certRepo := repository.NewCertificateRepository(db)
	wcRepo := repository.NewWorkerCompanyRepository(db)

	// Services
	companySvc := service.NewCompanyService(companyRepo)
	worksiteSvc := service.NewWorksiteService(worksiteRepo)
	workerSvc := service.NewWorkerService(workerRepo, certRepo, wcRepo)

	// Handlers
	companyHandler := handler.NewCompanyHandler(companySvc)
	worksiteHandler := handler.NewWorksiteHandler(worksiteSvc)
	workerHandler := handler.NewWorkerHandler(workerSvc)
	authHandler := handler.NewAuthHandler(authProvider)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.CORS(cfg.CORS.Origins))

	// Public routes
	r.Get("/health", handler.Health)
	r.Mount("/api/v1/auth", authHandler.Routes())

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authProvider))
		r.Mount("/api/v1/companies", companyHandler.Routes())
		r.Mount("/api/v1/worksites", worksiteHandler.Routes())
		r.Mount("/api/v1/workers", workerHandler.Routes())
	})

	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
