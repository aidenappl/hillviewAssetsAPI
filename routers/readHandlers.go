package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/query"
)

func ReadUserByTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["id"]

	user, err := query.ReadUser(db.DB, nil, &tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func ReadUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	userIntID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := query.ReadUser(db.DB, &userIntID, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func ReadAssetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["id"]

	intAssetID, err := strconv.Atoi(assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	asset, err := query.ReadAsset(db.DB, &intAssetID, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(asset)
}

func ReadAssetByTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["id"]

	asset, err := query.ReadAsset(db.DB, nil, &assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(asset)
}
