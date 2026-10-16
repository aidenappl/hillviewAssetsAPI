package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hillview.tv/assetsAPI/background"
	"github.com/hillview.tv/assetsAPI/env"
	"github.com/hillview.tv/assetsAPI/middleware"
	"github.com/hillview.tv/assetsAPI/routers"
	"github.com/hillview.tv/assetsAPI/telemetry"
)

func main() {
	telemetry.Init(telemetry.Config{
		Service:   "hillview-assets-api",
		IngestURL: env.MonitorIngestURL,
		APIKey:    env.MonitorAPIKey,
		Env:       env.MonitorEnv,
		SpoolDir:  env.MonitorSpoolDir,
		Stdout:    env.MonitorStdout,
	})

	primary := mux.NewRouter()

	// Healthcheck Endpoint

	primary.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Define the API Endpoints
	r := primary.PathPrefix("/assets/v1.1").Subrouter()

	// Add context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Monitor: one http.request.end event per request, panics recovered
	r.Use(telemetry.Middleware)

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// Adding response headers
	r.Use(middleware.MuxHeaderMiddleware)

	// == Checkout Routers ==

	r.HandleFunc("/checkout", routers.CheckoutHandler).Methods(http.MethodPost)
	r.HandleFunc("/checkin", routers.CheckinHandler).Methods(http.MethodPost)

	// == Create Routers ==

	create := r.PathPrefix("/create").Subrouter()

	create.HandleFunc("/asset", routers.CreateAssetHandler).Methods(http.MethodPost)
	create.HandleFunc("/user", routers.CreateUserHandler).Methods(http.MethodPost)

	// == Validator Routers ==

	valid := r.PathPrefix("/valid").Subrouter()

	valid.HandleFunc("/assetByID/{id}", routers.ValidAssetID).Methods(http.MethodGet)
	valid.HandleFunc("/assetByTag/{id}", routers.ValidAssetTag).Methods(http.MethodGet)
	valid.HandleFunc("/userByID/{id}", routers.ValidUserID).Methods(http.MethodGet)
	valid.HandleFunc("/userByTag/{id}", routers.ValidUserTag).Methods(http.MethodGet)

	// == Read Routers ==

	read := r.PathPrefix("/read").Subrouter()

	read.HandleFunc("/assetByID/{id}", routers.ReadAssetByIDHandler).Methods(http.MethodGet)
	read.HandleFunc("/assetByTag/{id}", routers.ReadAssetByTagHandler).Methods(http.MethodGet)
	read.HandleFunc("/userByID/{id}", routers.ReadUserByIDHandler).Methods(http.MethodGet)
	read.HandleFunc("/userByTag/{id}", routers.ReadUserByTagHandler).Methods(http.MethodGet)
	read.HandleFunc("/assetCheckoutHistory", routers.ReadAssetCheckoutHistory).Methods(http.MethodGet)

	// == Search Routers ==

	search := r.PathPrefix("/search").Subrouter()

	search.HandleFunc("/checkouts", routers.SearchCheckoutsHandler).Methods(http.MethodGet)

	// Launch API Listener
	fmt.Printf("✅ Hillview Assets API running on port %s\n", env.Port)

	// Add the health check polling
	telemetry.Go(ctx, "background.healthcheck", background.StartHealthCheckPolling)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"https://team.hillview.tv", "https://hillview.tv", "https://live.hillview.tv", "http://localhost:3000", "http://localhost:3001", "http://localhost:3002"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	exposedOk := handlers.ExposedHeaders([]string{"X-Request-ID"})

	server := &http.Server{
		Addr:              ":" + env.Port,
		Handler:           handlers.CORS(originsOk, headersOk, exposedOk, methodsOk)(primary),
		ReadHeaderTimeout: 15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	if err := telemetry.ServeWithGracefulShutdown(server, stop); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
