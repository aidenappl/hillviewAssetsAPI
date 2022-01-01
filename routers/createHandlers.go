package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/query"
)

type CreateAssetRequest struct {
	Name        *string `json:"name"`
	ImageURL    *string `json:"image_url"`
	Identifier  *string `json:"identifier"`
	Description *string `json:"description"`
	Category    *int    `json:"category"`

	SerialNumber *string `json:"serial_number"`
	Manufacturer *string `json:"manufacturer"`
	Model        *string `json:"model"`
	Notes        *string `json:"notes"`
}

func CreateAssetHandler(w http.ResponseWriter, r *http.Request) {
	body := CreateAssetRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Name == nil {
		http.Error(w, "body must include name", http.StatusBadRequest)
		return
	}

	if body.Identifier == nil {
		http.Error(w, "body must include an identifier", http.StatusBadRequest)
		return
	}

	if body.Category == nil {
		http.Error(w, "body must include a category", http.StatusBadRequest)
		return
	}

	assetID, err := query.CreateNewAsset(db.DB, *body.Name, *body.Identifier, *body.Category, body.ImageURL, body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if body.SerialNumber != nil || body.Manufacturer != nil || body.Model != nil || body.Notes != nil {
		err = query.CreateNewAssetInfo(db.DB, assetID, body.SerialNumber, body.Manufacturer, body.Model, body.Notes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

}
