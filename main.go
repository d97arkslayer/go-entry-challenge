package main

import (
	"github.com/d97arkslayer/go-entry-challenge/Router"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
)
/**
 * init
 * Use to init the env variables
 */
func init() {
	// Load Env vars
	gotenv.Load()
}

/**
 * main
 * Setup server and serve
 */
func main()  {
	// Get api routes
	router := Router.Router()
	// Get Port to listen and serve
	PORT := os.Getenv("PORT")
	// Check port variable
	if PORT == "" {
		// Set port variable
		PORT = "5000"
	}
	// Serve the api
	log.Fatal(http.ListenAndServe(":" + PORT, router))
}