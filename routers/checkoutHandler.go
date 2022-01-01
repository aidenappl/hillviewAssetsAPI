package routers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/query"
)

type CheckoutRequest struct {
	UserID   int       `json:"user_id"`
	AssetID  int       `json:"asset_id"`
	Duration time.Time `json:"duration"`
	Notes    *string   `json:"notes"`
	OffSite  bool      `json:"offsite"`
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	body := CheckoutRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	activeCheckout, err := query.ReadActiveCheckouts(db.DB, body.AssetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if activeCheckout != nil {
		http.Error(w, "Asset is already checked out", http.StatusBadRequest)
		return
	}

	err = query.CheckoutAsset(db.DB, body.UserID, body.AssetID, body.Duration, body.OffSite, body.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type CheckinRequest struct {
	AssetID int     `json:"asset_id"`
	Notes   *string `json:"notes"`
}

func CheckinHandler(w http.ResponseWriter, r *http.Request) {
	body := CheckinRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	activeCheckout, err := query.ReadActiveCheckouts(db.DB, body.AssetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if activeCheckout == nil {
		http.Error(w, "Asset is already checked in", http.StatusBadRequest)
		return
	}

	err = query.CheckinAsset(db.DB, body.AssetID, body.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
