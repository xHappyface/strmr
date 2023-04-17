package handlers

import (
	"net/http"
	"text/template"
)

func (h *Handlers) YouTubeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("./templates/youtube.html"))
		tmpl.Execute(w, struct {
			Title      string
			Javascript []string
			CSS        []string
		}{
			Title: "OBS stream settings",
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"obs",
			},
			CSS: []string{
				"obs",
			},
		})
	}
}
