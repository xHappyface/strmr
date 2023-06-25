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
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
}

func (h *Handlers) TwitchUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data TwitchUpdate
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.Unmarshal(reqBody, &data)
		err = h.twitch.ChangeStream("jnrprgmr", data.Title, data.CategoryID, data.Tags)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t, err := h.database.GetLatestMetadataByKey("title", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(t) != 1 || t[0].MetadataValue != data.Title {
			err = h.database.InsertMetadata("title", data.Title)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		d, err := h.database.GetLatestMetadataByKey("description", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(d) != 1 || d[0].MetadataValue != data.Description {
			err = h.database.InsertMetadata("description", data.Description)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		c, err := h.database.GetLatestMetadataByKey("category", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(c) != 1 || c[0].MetadataValue != data.CategoryName {
			err = h.database.InsertMetadata("category", data.CategoryName)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		ca, err := h.database.GetCategoryByName(data.CategoryName)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ca == nil {
			err = h.database.InsertCategory(data.CategoryName)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		ts, err := h.database.GetLatestMetadataByKey("tags", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tags_string := strings.Join(data.Tags, ",")
		if len(ts) != 1 || t[0].MetadataValue != tags_string {
			err = h.database.InsertMetadata("tags", tags_string)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
