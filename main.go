package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hillview.tv/assetsAPI/env"
	"github.com/hillview.tv/assetsAPI/middleware"
)

func main() {
	r := mux.NewRouter()

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// Adding response headers
	r.Use(middleware.MuxHeaderMiddleware)

	// Launch API Listener
	fmt.Printf("✅ Hillview Assets API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
