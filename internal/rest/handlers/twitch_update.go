package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TwitchUpdate struct {
	Title string `json:"title"`
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
		err = h.twitch.ChangeStreamTitle("jnrprgmr", data.Title)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
