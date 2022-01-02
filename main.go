package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hillview.tv/assetsAPI/env"
	"github.com/hillview.tv/assetsAPI/middleware"
	"github.com/hillview.tv/assetsAPI/routers"
)

func main() {
	primary := mux.NewRouter()

	// Healthcheck Endpoint

	primary.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	}).Methods(http.MethodGet)

	r := primary.PathPrefix("/assets/v1.1").Subrouter()

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// Adding response headers
	r.Use(middleware.MuxHeaderMiddleware)

	// == Primary Checkout API ==

	r.HandleFunc("/checkout", routers.CheckoutHandler).Methods(http.MethodPost)
	r.HandleFunc("/checkin", routers.CheckinHandler).Methods(http.MethodPost)

	// == Create Routers ==

	create := r.PathPrefix("/create").Subrouter()

	create.HandleFunc("/asset", routers.CreateAssetHandler).Methods(http.MethodPost)

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

	// Launch API Listener
	fmt.Printf("✅ Hillview Assets API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With, Content-Type, Origin, Authorization, Accept, X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(primary)))
}
