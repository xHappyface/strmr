package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jnrprgmr/dog/pkg/obs"
)

func (h *Handlers) UpdateOBSTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data obs.Task
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		status := http.StatusOK
		err = h.obs.SetTask(data)
		if err != nil {
			status = http.StatusBadRequest

		}
		w.WriteHeader(status)
	}
}
