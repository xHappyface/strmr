package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/nicklaw5/helix"
)

func (h *Handlers) TwitchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userAccessToken := h.twitch.Client.GetUserAccessToken()
		authorized, _, err := h.twitch.Client.ValidateToken(userAccessToken)
		if err != nil {
			// handle error
		}
		url := h.twitch.Client.GetAuthorizationURL(&helix.AuthorizationURLParams{
			ResponseType: "code",
			Scopes:       []string{"channel:manage:broadcast"},
			State:        "some-statedasdad",
			ForceVerify:  false,
		})
		if !authorized {
			if h.refreshToken != "" {
				resp, err := h.twitch.Client.RefreshUserAccessToken(h.refreshToken)
				if err != nil {
					fmt.Println("coult not refresh access token")
				} else {
					h.token = resp.Data.AccessToken
					h.refreshToken = resp.Data.RefreshToken
					h.twitch.Client.SetUserAccessToken(h.token)
					userAccessToken := h.twitch.Client.GetUserAccessToken()
					authorized, _, _ = h.twitch.Client.ValidateToken(userAccessToken)
				}
			}
		}
		tmpl := template.Must(template.ParseFiles("./templates/twitch.html"))
		tmpl.Execute(w, struct {
			Title      string
			Authorized bool
			AuthURL    string
			Javascript []string
			CSS        []string
		}{
			Title:      "Twitch stream settings",
			Authorized: authorized,
			AuthURL:    url,
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"twitch",
			},
			CSS: []string{
				"test",
			},
		})
	}
}
