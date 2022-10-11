package routers

import (
	"net/http"
)

func SearchCheckoutsHandler(w http.ResponseWriter, r *http.Request) {
	// userID, _ := r.URL.Query()["userID"]
	// tag, _ := r.URL.Query()["tag"]

	// request := query.SearchCheckoutsRequest{}

	// if len(userID) > 0 {
	// 	request.UserID = &userID[0]
	// }

	// if len(tag) > 0 {
	// 	request.Tag = &tag[0]
	// }

	// history, err := query.SearchCheckouts(db.DB, request)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(history)
}
