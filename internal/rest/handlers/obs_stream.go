package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
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
			err = h.database.EndActiveStreams()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			if data.Stream {
				err = h.database.InsertStream()
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			_, err := h.obs.ToggleStream()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		record_status, err := h.obs.GetRecordStatus()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if data.Record != record_status {
			err = h.database.EndActiveMediaRecordings()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if data.Record {
				dir, err := h.obs.GetRecordDirectory()
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = h.obs.SetProfileParameter("FilenameFormatting", time.Now().UTC().Format(time.RFC3339))
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
					return
				}
				file_name, err := h.obs.GetProfileParameter("FilenameFormatting")
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = h.database.InsertMediaRecording(file_name, dir)
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			err := h.obs.ToggleRecord()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
}
