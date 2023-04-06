package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Stream struct {
	Stream     bool    `json:"stream"`
	Record     bool    `json:"record"`
	OutputFile *string `json:"output_file,omitempty"`
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
	}
}
