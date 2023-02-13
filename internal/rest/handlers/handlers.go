package handlers

import "github.com/jnrprgmr/dog/pkg/twitch"

type Handlers struct {
	twitch       *twitch.Twitch
	code         string
	token        string
	refreshToken string
}

func New(twitchCli *twitch.Twitch) *Handlers {
	return &Handlers{
		twitch:       twitchCli,
		code:         "",
		token:        "",
		refreshToken: "",
	}
}
