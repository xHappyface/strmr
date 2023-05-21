package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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
		if len(data.Text) != 0 {
			err = h.database.InsertMetadata("task", data.Text)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			text_color, err := obs.ConvertColor(data.Color)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			config := strconv.Itoa(int(text_color)) + "," + fmt.Sprintf("%f", data.Width) + "," + fmt.Sprintf("%f", data.Height) + "," + fmt.Sprintf("%f", data.PosX) + "," + fmt.Sprintf("%f", data.PosY)
			err = h.database.InsertMetadata("task_config", config)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if data.Background != nil {
			background_color, err := obs.ConvertColor(data.Background.Color)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			config := strconv.Itoa(int(background_color)) + "," + fmt.Sprintf("%f", data.Width) + "," + fmt.Sprintf("%f", data.Height) + "," + fmt.Sprintf("%f", data.PosX) + "," + fmt.Sprintf("%f", data.PosY)
			err = h.database.InsertMetadata("task_background_config", config)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		err = h.obs.SetTask(data)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
