package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"text/template"
	"time"
)

var talking = false

type avatarRequest struct {
	Text string `json:"text"`
}

func (h *Handlers) AvatarStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data avatarRequest
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		talking = true
		cmd := exec.Command("espeak", "-x", data.Text)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		end := time.Now()
		t := end.Sub(start)
		fmt.Println("ere")
		err = h.database.InsertSubtitle(data.Text, t.Seconds())
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//fmt.Println(t.Seconds())
		talking = false
	} else if r.Method == http.MethodGet {
		if talking {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (h *Handlers) Avatar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./templates/avatar.html"))
		tmpl.Execute(w, struct {
			Title      string
			Javascript []string
			CSS        []string
		}{
			Title: "OBS avatar widget",
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"avatar",
			},
			CSS: []string{
				"avatar",
			},
		})
	}
}
