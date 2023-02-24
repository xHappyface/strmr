package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CreateSceneReq struct {
	Name string `json:"name"`
}

type SceneResp struct {
	Names []string `json:"names"`
}

func (h *Handlers) CreateScene(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data CreateSceneReq
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		scenes, err := h.obs.GetSceneList()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(data.Name) == 0 {
			h.ErrorResponse(w, "input name must be greater than length 0", http.StatusBadRequest)
			return
		}
		names := []string{}
		for i := range scenes.Scenes {
			scene := scenes.Scenes[i].SceneName
			names = append(names, scene)
			if scene == data.Name {
				h.ErrorResponse(w, "Scene Already Exists", http.StatusBadRequest)
				return
			}
		}
		_, err = h.obs.CreateScene(data.Name)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		names = append(names, data.Name)
		resp := SceneResp{
			Names: names,
		}
		b, err := json.Marshal(resp)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
