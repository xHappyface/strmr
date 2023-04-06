package handlers

import (
	"net/http"
	"text/template"

	"github.com/andreykaipov/goobs/api/typedefs"
)

type Task struct {
	Text            string
	Color           string
	PosX            float64
	PosY            float64
	Width           float64
	Height          float64
	Background      bool
	BackgroundColor string
}

func (h *Handlers) ObsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scenes, err := h.obs.GetSceneList()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		input_settings, err := h.obs.GetInputSettings("test")
		input_exists := true
		if err != nil {
			input_exists = false
		}
		background_settings, err := h.obs.GetInputSettings("background")
		background_exists := true
		if err != nil {
			background_exists = false
		}
		stream_status, err := h.obs.GetStreamStatus()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		record_status, err := h.obs.GetRecordStatus()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		task := Task{
			Text:       "",
			Color:      "000000",
			PosX:       0,
			PosY:       0,
			Width:      0,
			Height:     0,
			Background: background_exists,
		}
		if input_exists {
			transform, err := h.obs.GetSceneItemTransform(h.obs.GetSceneItemId("Main", "test"), "Main")
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			hex, err := h.obs.ConvertIntToHex(int64(input_settings.InputSettings["color1"].(float64)))
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			hx := *hex
			bHex := hx[2:4]
			gHex := hx[4:6]
			rHex := hx[6:8]
			color := rHex + gHex + bHex
			task = Task{
				Text:       input_settings.InputSettings["text"].(string),
				Color:      color,
				PosX:       transform.SceneItemTransform.PositionX,
				PosY:       transform.SceneItemTransform.PositionY,
				Width:      transform.SceneItemTransform.BoundsWidth,
				Height:     transform.SceneItemTransform.BoundsHeight,
				Background: background_exists,
			}
			if background_exists {
				hex, err := h.obs.ConvertIntToHex(int64(background_settings.InputSettings["color"].(float64)))
				if err != nil {
					h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
				hx := *hex
				bHex := hx[2:4]
				gHex := hx[4:6]
				rHex := hx[6:8]
				task.BackgroundColor = rHex + gHex + bHex
			}
		}
		tmpl := template.Must(template.ParseFiles("./templates/obs.html"))
		tmpl.Execute(w, struct {
			Title        string
			Javascript   []string
			CSS          []string
			Scenes       []*typedefs.Scene
			Task         Task
			StreamStatus bool
			RecordStatus bool
		}{
			Title: "OBS stream settings",
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"obs",
			},
			CSS: []string{
				"obs",
			},
			Scenes:       scenes.Scenes,
			Task:         task,
			StreamStatus: stream_status,
			RecordStatus: record_status,
		})
	}
}
