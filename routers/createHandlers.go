package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/query"
)

type CreateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Tag      *string `json:"tag"`
	PhotoURL *string `json:"photo_url"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if body.Name == nil {
		http.Error(w, "Missing name", http.StatusBadRequest)
		return
	}

	if body.Email == nil {
		http.Error(w, "Missing email", http.StatusBadRequest)
		return
	}

	if body.Tag == nil {
		http.Error(w, "Missing tag", http.StatusBadRequest)
		return
	}

	if body.PhotoURL == nil {
		http.Error(w, "Missing photo_url", http.StatusBadRequest)
		return
	}

	err = query.CreateNewUser(db.DB, *body.Name, *body.Email, *body.Tag, *body.PhotoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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
