package Router

import (
	"github.com/d97arkslayer/go-entry-challenge/Controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
)

/**
 * Router
 * Use to init the CHI router, and add server routes
 */
func Router() *chi.Mux {
	router := chi.NewRouter()
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	// Routes
	router.Get("/", Controllers.IndexBuyers)

	return router
}