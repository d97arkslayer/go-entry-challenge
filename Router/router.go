package Router

import (
	"github.com/d97arkslayer/go-entry-challenge/Controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

/**
 * Router
 * Use to init the CHI router, and add server routes
 */
func Router(){
	router := chi.NewRouter()
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Get("/", Controllers.IndexBuyers)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	log.Fatal(http.ListenAndServe(":" + PORT, router))
}