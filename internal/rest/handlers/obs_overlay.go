package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jnrprgmr/strmr/pkg/obs"
)

type Overlay struct {
	Text            string    `json:"text"`
	TextWidth       float64   `json:"text_width"`
	TextHeight      float64   `json:"text_height"`
	TextPosX        float64   `json:"text_posx"`
	TextPosY        float64   `json:"text_posy"`
	TextColor       obs.Color `json:"text_color"`
	BackgroundColor obs.Color `json:"background_color"`
	Enabled         bool      `json:"enabled"`
}

func (h *Handlers) UpdateOBSOverlay(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data Overlay
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		current_scene, err := h.obs.GetCurrentScene()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		transform, err := h.obs.GetSceneItemTransform(h.obs.GetSceneItemId(current_scene, h.obs.ScreenSourceName), current_scene)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		color, err := obs.ConvertColor(data.BackgroundColor)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		overlay_background_settings := map[string]interface{}{
			"color":  color,
			"width":  transform.SceneItemTransform.Width,
			"height": transform.SceneItemTransform.Height,
		}
		_, err = h.obs.SetInputSettings(h.obs.OverlayBackgroundSourceName, overlay_background_settings)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		text_color, err := obs.ConvertColor(data.TextColor)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		overlay_text_settings := map[string]interface{}{
			"text":   data.Text,
			"color1": text_color,
			"color2": text_color,
		}
		_, err = h.obs.SetInputSettings(h.obs.OverlayTextSourceName, overlay_text_settings)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.obs.SetSceneItemTransform(h.obs.GetSceneItemId(current_scene, h.obs.OverlayTextSourceName), current_scene, data.TextPosX, data.TextPosY, data.TextWidth, data.TextHeight)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.obs.SetSceneItemEnabled(h.obs.GetSceneItemId(current_scene, h.obs.OverlayBackgroundSourceName), current_scene, data.Enabled)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.obs.SetSceneItemEnabled(h.obs.GetSceneItemId(current_scene, h.obs.OverlayTextSourceName), current_scene, data.Enabled)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
