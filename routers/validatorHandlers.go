package routers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/query"
)

func ValidUserTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["id"]

	user, err := query.ReadUser(db.DB, nil, &tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func ValidUserID(w http.ResponseWriter, r *http.Request) {
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

	if user != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func ValidAssetID(w http.ResponseWriter, r *http.Request) {
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

	if asset != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func ValidAssetTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["id"]

	asset, err := query.ReadAsset(db.DB, nil, &assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if asset != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
