package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TwitchSearchCategories struct {
	Query string `json:"query"`
}

func (h *Handlers) TwitchSearchCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data TwitchSearchCategories
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		categories, err := h.twitch.SearchCategories(data.Query)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(categories)
	}
}
