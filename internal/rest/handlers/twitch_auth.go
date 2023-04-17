package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) TwitchAuthHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	h.twitch.Code = code
	resp, err := h.twitch.Client.RequestUserAccessToken(h.twitch.Code)
	if err != nil {
		// handle error
		fmt.Println("test")
	}
	// Set the access token on the client
	h.twitch.Token = resp.Data.AccessToken
	h.twitch.RefreshToken = resp.Data.RefreshToken
	h.twitch.Client.SetUserAccessToken(h.twitch.Token)
	users, err := h.twitch.GetUsers([]string{})
	if err != nil {
		users = map[string]string{}
	}
	if id, ok := users["jnrprgmr"]; ok {
		err = h.database.InsertUser("twitch", id)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.database.InsertMetadata("refresh_token", h.twitch.RefreshToken)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.database.InsertMetadata("access_token", h.twitch.Token)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "http://localhost:8080/twitch", http.StatusSeeOther)
}
