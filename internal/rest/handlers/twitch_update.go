package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type TwitchUpdate struct {
	Title        string   `json:"title"`
	CategoryID   string   `json:"category_id"`
	CategoryName string   `json:"category_name"`
	Tags         []string `json:"tags"`
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
		err = h.twitch.ChangeStream("jnrprgmr", data.Title, data.CategoryID, data.Tags)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t, err := h.database.GetLatestMetadataByKey("title", 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(t) != 1 {
			h.ErrorResponse(w, "Error getting latest title from metadata", http.StatusInternalServerError)
			return
		}
		title := t[0]
		if title.MetadataValue != data.Title {
			err = h.database.InsertMetadata("title", data.Title)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		c, err := h.database.GetLatestMetadataByKey("category", 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(c) != 1 {
			h.ErrorResponse(w, "Error getting latest category from metadata", http.StatusInternalServerError)
			return
		}
		category := c[0]
		if category.MetadataValue != data.CategoryName {
			err = h.database.InsertMetadata("category", data.CategoryName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		ca, err := h.database.GetCategoryByName(data.CategoryName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if ca == nil {
			err = h.database.InsertCategory(data.CategoryName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		ts, err := h.database.GetLatestMetadataByKey("tags", 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(ts) != 1 {
			h.ErrorResponse(w, "Error getting latest tags from metadata", http.StatusInternalServerError)
			return
		}
		tags := ts[0]
		tags_string := strings.Join(data.Tags, ",")
		if tags.MetadataValue != tags_string {
			err = h.database.InsertMetadata("tags", tags_string)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
