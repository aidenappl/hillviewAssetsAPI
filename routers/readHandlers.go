package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/query"
)

func ReadAssetCheckoutHistory(w http.ResponseWriter, r *http.Request) {
	id, _ := r.URL.Query()["id"]
	tag, _ := r.URL.Query()["tag"]

	request := query.ReadAssetCheckoutsRequest{
		Limit: 10,
	}

	if len(id) > 0 {
		request.AssetID = &id[0]
	}

	if len(tag) > 0 {
		request.TagID = &tag[0]
	}

	fmt.Println(request)

	history, err := query.ReadAssetCheckoutHistory(db.DB, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

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
