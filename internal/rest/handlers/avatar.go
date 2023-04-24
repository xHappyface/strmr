package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"text/template"
	"time"
)

var talking = false

func (h *Handlers) AvatarStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		talking = true
		s := "Make the computer speak for a long time without taking a break and get the time amount"
		cmd := exec.Command("espeak", "-x", s)
		start := time.Now()
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		end := time.Now()
		t := end.Sub(start)
		fmt.Println(t)
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
