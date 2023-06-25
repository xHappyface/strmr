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
		access_token, err := h.database.GetLatestMetadataByKey("access_token", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(access_token) == 1 {
			h.twitch.Token = access_token[0].MetadataValue
			h.twitch.Client.SetUserAccessToken(h.twitch.Token)
		}
		refresh_token, err := h.database.GetLatestMetadataByKey("refresh_token", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(refresh_token) == 1 {
			h.twitch.RefreshToken = refresh_token[0].MetadataValue
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
		description := ""
		descriptions_hist, err := h.database.GetLatestMetadataByKey("description", 1)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(descriptions_hist) != 0 {
			description = descriptions_hist[0].MetadataValue
		}
		game_titles := []string{}
		categories, err := h.database.GetDistinctMetadataValuesByKey("category")
		if err != nil {
			game_titles = []string{}
		}
		for k := range categories {
			game_titles = append(game_titles, categories[k].MetadataValue)
		}
		games, err := h.twitch.GetGames(game_titles)
		if err != nil {
			games = map[string]string{}
		}
		title_hist, err := h.database.GetLatestMetadataByKey("title", 5)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		titles := []string{}
		for i := range title_hist {
			titles = append(titles, title_hist[i].MetadataValue)
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
			Title       string
			Authorized  bool
			AuthURL     string
			Description string
			Games       map[string]GameElement
			Users       map[string]string
			Channel     twitch.Channel
			Titles      []string
			Javascript  []string
			CSS         []string
		}{
			Title:       "Twitch stream settings",
			Authorized:  authorized,
			Games:       g,
			Users:       users,
			Channel:     channel,
			Description: description,
			Titles:      titles,
			AuthURL:     url,
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
