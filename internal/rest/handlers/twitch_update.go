package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TwitchUpdate struct {
	Title        string `json:"title"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

func (h *Handlers) TwitchUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data TwitchUpdate
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		err = h.twitch.ChangeStream("jnrprgmr", data.Title, data.CategoryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		title, err := h.database.GetLatestMetadataByKey("title")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if title.MetadataValue != data.Title {
			err = h.database.InsertMetadata("title", data.Title)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		category, err := h.database.GetLatestMetadataByKey("category")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if category.MetadataValue != data.CategoryName {
			err = h.database.InsertMetadata("category", data.CategoryName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
