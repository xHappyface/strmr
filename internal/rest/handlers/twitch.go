package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/nicklaw5/helix/v2"
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
			if h.twitch.RefreshToken != "" {
				resp, err := h.twitch.Client.RefreshUserAccessToken(h.twitch.RefreshToken)
				if err != nil {
					fmt.Println("coult not refresh access token")
				} else {
					h.twitch.Token = resp.Data.AccessToken
					h.twitch.RefreshToken = resp.Data.RefreshToken
					h.twitch.Client.SetUserAccessToken(h.twitch.Token)
					userAccessToken := h.twitch.Client.GetUserAccessToken()
					authorized, _, _ = h.twitch.Client.ValidateToken(userAccessToken)
				}
			}
		}
		users, err := h.twitch.GetUsers([]string{})
		if err != nil {
			users = map[string]string{}
		}
		if id, ok := users["jnrprgmr"]; ok {
			err = h.database.InsertUser("twitch", id)
			if err != nil {
				fmt.Println("couldnt insert user in db from handler")
			}
		}
		games, err := h.twitch.GetGames([]string{"Dota 2", "Software and Game Development", "pokemon"})
		if err != nil {
			games = map[string]string{}
		}
		tmpl := template.Must(template.ParseFiles("./templates/twitch.html"))
		tmpl.Execute(w, struct {
			Title      string
			Authorized bool
			AuthURL    string
			Games      map[string]string
			Users      map[string]string
			Javascript []string
			CSS        []string
		}{
			Title:      "Twitch stream settings",
			Authorized: authorized,
			Games:      games,
			Users:      users,
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
