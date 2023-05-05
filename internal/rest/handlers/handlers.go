package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jnrprgmr/strmr/pkg/database"
	"github.com/jnrprgmr/strmr/pkg/obs"
	"github.com/jnrprgmr/strmr/pkg/twitch"
	"github.com/jnrprgmr/strmr/pkg/youtube"
)

type Handlers struct {
	twitch   *twitch.Twitch
	obs      *obs.OBS
	youtube  *youtube.YouTube
	database *database.Database
}

type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (h *Handlers) ErrorResponse(w http.ResponseWriter, message string, code int) {
	err_resp := HTTPError{
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(err_resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(b)
}

func New(twitchCli *twitch.Twitch, obsCli *obs.OBS, yt *youtube.YouTube, db *database.Database) *Handlers {
	return &Handlers{
		twitch:   twitchCli,
		obs:      obsCli,
		youtube:  yt,
		database: db,
	}
}
