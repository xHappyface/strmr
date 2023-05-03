package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/jnrprgmr/strmr/pkg/twitch"
	"github.com/nicklaw5/helix/v2"
)

type GameElement struct {
	ID       string
	Selected bool
}

func (h *Handlers) TwitchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		access_token, err := h.database.GetLatestMetadataByKey("access_token")
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if access_token != nil {
			h.twitch.Token = access_token.MetadataValue
			h.twitch.Client.SetUserAccessToken(h.twitch.Token)
		}
		refresh_token, err := h.database.GetLatestMetadataByKey("refresh_token")
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if refresh_token != nil {
			h.twitch.RefreshToken = refresh_token.MetadataValue
		}
		userAccessToken := h.twitch.Client.GetUserAccessToken()
		authorized, _, err := h.twitch.Client.ValidateToken(userAccessToken)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url := h.twitch.Client.GetAuthorizationURL(&helix.AuthorizationURLParams{
			ResponseType: "code",
			Scopes:       []string{"channel:manage:broadcast"},
			State:        "some-statedasd",
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
					h.twitch.Client.SetUserAccessToken(h.twitch.Token)
					userAccessToken := h.twitch.Client.GetUserAccessToken()
					authorized, _, _ = h.twitch.Client.ValidateToken(userAccessToken)
				}
			}
		}
		users, err := h.twitch.GetUsers([]string{})
		if err != nil || len(users) != 1 {
			users = map[string]string{}
		}
		user_login := ""
		user_id := ""
		for k := range users {
			user_login = k
			user_id = users[k]
		}
		titles := []string{}
		categories, err := h.database.GetDistinctMetadataValuesByKey("category")
		if err != nil {
			titles = []string{}
		}
		for k := range categories {
			titles = append(titles, categories[k].MetadataValue)
		}
		games, err := h.twitch.GetGames(titles)
		if err != nil {
			games = map[string]string{}
		}
		channel := twitch.Channel{}
		channels, _ := h.twitch.GetChannelInformation([]string{user_id})
		if ch, ok := channels[user_login]; ok {
			channel = ch
		}
		g := map[string]GameElement{}
		for k, v := range games {
			s := false
			if channel.CategoryID == v {
				s = true
			}
			ga := GameElement{
				ID:       v,
				Selected: s,
			}
			g[k] = ga
		}
		tmpl := template.Must(template.ParseFiles("./templates/twitch.html"))
		tmpl.Execute(w, struct {
			Title      string
			Authorized bool
			AuthURL    string
			Games      map[string]GameElement
			Users      map[string]string
			Channel    twitch.Channel
			Javascript []string
			CSS        []string
		}{
			Title:      "Twitch stream settings",
			Authorized: authorized,
			Games:      g,
			Users:      users,
			Channel:    channel,
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
