package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type YouTubeCategory struct {
	CategoryName string `json:"category_name"`
	RelatedID    string `json:"related_id"`
}

func (h *Handlers) YouTubeCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data YouTubeCategory
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		if data.CategoryName == "" || data.RelatedID == "" {
			h.ErrorResponse(w, "Cannot use empty values when updating category", http.StatusBadRequest)
			return
		}
		cat, err := h.database.GetCategoryByName(data.CategoryName)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cat == nil {
			h.ErrorResponse(w, "Attempting to correlate category that has never been used", http.StatusBadRequest)
			return
		}
		err = h.database.UpdateCategoryByName(data.RelatedID, data.CategoryName)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
