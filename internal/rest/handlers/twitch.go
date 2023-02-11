package handlers

import (
	"fmt"
	"net/http"

	"github.com/nicklaw5/helix"
)

func (h *Handlers) TwitchHandler(w http.ResponseWriter, r *http.Request) {

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
	auth_msg := "Not authorized"
	auth_color := "red"
	if authorized {
		auth_msg = "Authorized"
		auth_color = "green"
		h.twitch.ChangeStreamTitle("jnrprgmr", "Change title")
	}
	fmt.Fprintf(w, `<h1>Settings: <span style='color:%s'>%s</span></h1><br><a href=%s>Authenticate</a>`, auth_color, auth_msg, url)
}
