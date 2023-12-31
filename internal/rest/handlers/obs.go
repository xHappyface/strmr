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

type FlatOverlay struct {
	Text            string
	TextWidth       float64
	TextHeight      float64
	TextPosX        float64
	TextPosY        float64
	TextColor       string
	BackgroundColor string
	Enabled         bool
}

func (h *Handlers) ObsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scenes, err := h.obs.GetSceneList()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		input_settings, err := h.obs.GetInputSettings(h.obs.TaskSourceName)
		input_exists := true
		if err != nil {
			input_exists = false
		}
		background_settings, err := h.obs.GetInputSettings(h.obs.BackgroundSourceName)
		background_exists := true
		if err != nil {
			background_exists = false
		}
		overlay_text_settings, err := h.obs.GetInputSettings(h.obs.OverlayTextSourceName)
		overlay_text_exists := true
		if err != nil {
			overlay_text_exists = false
		}
		overlay_background_settings, err := h.obs.GetInputSettings(h.obs.OverlayBackgroundSourceName)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
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
			current_scene, err := h.obs.GetCurrentScene()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			transform, err := h.obs.GetSceneItemTransform(h.obs.GetSceneItemId(current_scene, h.obs.TaskSourceName), current_scene)
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
		overlay := FlatOverlay{
			Text:            "",
			TextWidth:       1,
			TextHeight:      1,
			TextPosX:        0,
			TextPosY:        0,
			TextColor:       "000000",
			BackgroundColor: "000000",
			Enabled:         false,
		}
		if overlay_text_exists {
			current_scene, err := h.obs.GetCurrentScene()
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			overlay_text_transform, err := h.obs.GetSceneItemTransform(h.obs.GetSceneItemId(current_scene, h.obs.OverlayTextSourceName), current_scene)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			overlay_text_hex, err := h.obs.ConvertIntToHex(int64(overlay_text_settings.InputSettings["color1"].(float64)))
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			overlay_background_hex, err := h.obs.ConvertIntToHex(int64(overlay_background_settings.InputSettings["color"].(float64)))
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			overlay.Text = overlay_text_settings.InputSettings["text"].(string)
			hx := *overlay_text_hex
			bHex := hx[2:4]
			gHex := hx[4:6]
			rHex := hx[6:8]
			color := rHex + gHex + bHex
			overlay.TextColor = color
			overlay.TextPosX = overlay_text_transform.SceneItemTransform.PositionX
			overlay.TextPosY = overlay_text_transform.SceneItemTransform.PositionY
			overlay.TextWidth = overlay_text_transform.SceneItemTransform.BoundsWidth
			overlay.TextHeight = overlay_text_transform.SceneItemTransform.BoundsHeight
			bghx := *overlay_background_hex
			bgbHex := bghx[2:4]
			bggHex := bghx[4:6]
			bgrHex := bghx[6:8]
			bgcolor := bgrHex + bggHex + bgbHex
			overlay.BackgroundColor = bgcolor
			visible, err := h.obs.GetSceneSourceVisible(h.obs.GetSceneItemId(current_scene, h.obs.OverlayTextSourceName), current_scene)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			overlay.Enabled = *visible
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
			Overlay      FlatOverlay
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
			Overlay:      overlay,
		})
	}
}
