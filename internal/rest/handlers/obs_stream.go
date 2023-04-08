package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Stream struct {
	Stream bool `json:"stream"`
	Record bool `json:"record"`
}

func (h *Handlers) UpdateOBSStream(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data Stream
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		stream_status, err := h.obs.GetStreamStatus()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if data.Stream != stream_status {
			_, err := h.obs.ToggleStream()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			if data.Stream {
				err = h.database.EndActiveStreams()
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
				err = h.database.InsertStream()
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
		record_status, err := h.obs.GetRecordStatus()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if data.Record != record_status {
			err := h.obs.ToggleRecord()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
}
