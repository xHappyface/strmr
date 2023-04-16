package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jnrprgmr/strmr/pkg/obs"
)

func (h *Handlers) UpdateOBSTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data obs.Task
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		err = h.database.InsertTask(data.Text)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.obs.SetTask(data)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
