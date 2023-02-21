package handlers

import (
	"github.com/jnrprgmr/dog/pkg/database"
	"github.com/jnrprgmr/dog/pkg/obs"
	"github.com/jnrprgmr/dog/pkg/twitch"
)

type Handlers struct {
	twitch   *twitch.Twitch
	obs      *obs.OBS
	database *database.Database
}

func New(twitchCli *twitch.Twitch, obsCli *obs.OBS, db *database.Database) *Handlers {
	return &Handlers{
		twitch:   twitchCli,
		obs:      obsCli,
		database: db,
	}
}
