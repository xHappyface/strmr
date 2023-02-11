package handlers

import (
	"net/http"
)

func (h *Handlers) TwitchAuthHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	h.code = code
	resp, err := h.twitch.Client.RequestUserAccessToken(h.code)
	if err != nil {
		// handle error
	}

	// Set the access token on the client
	h.token = resp.Data.AccessToken
	h.twitch.Client.SetUserAccessToken(h.token)
	http.Redirect(w, r, "http://localhost:8080/twitch", http.StatusSeeOther)
}
