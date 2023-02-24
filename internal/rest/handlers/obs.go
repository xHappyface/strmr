package handlers

import (
	"net/http"
	"text/template"

	"github.com/andreykaipov/goobs/api/typedefs"
)

func (h *Handlers) ObsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scenes, err := h.obs.GetSceneList()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tmpl := template.Must(template.ParseFiles("./templates/obs.html"))
		tmpl.Execute(w, struct {
			Title      string
			Javascript []string
			CSS        []string
			Scenes     []*typedefs.Scene
		}{
			Title: "OBS stream settings",
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"obs",
			},
			CSS: []string{
				"test",
			},
			Scenes: scenes.Scenes,
		})
	}
}
