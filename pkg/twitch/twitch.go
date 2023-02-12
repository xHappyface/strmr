package twitch

import (
	"errors"

	"github.com/nicklaw5/helix"
)

type Twitch struct {
	Client *helix.Client
}

func New(client *helix.Client) *Twitch {
	return &Twitch{
		Client: client,
	}
}

func (t *Twitch) refreshAppToken() error {
	auth, err := t.Client.RequestAppAccessToken([]string{})
	if err != nil || auth.StatusCode != 200 {
		return errors.New("Error generating app token from twitch: " + err.Error())
	}
	t.Client.SetAppAccessToken(auth.Data.AccessToken)
	return nil
}

func (t *Twitch) ChangeStreamTitle(username string, title string) error {
	err := t.refreshAppToken()
	if err != nil {
		return errors.New("Error refreshing token: " + err.Error())
	}
	users, err := t.GetUser([]string{username, "934209qadijoiwjdolasdlaksdjlkj"})
	broadcaster_id, ok := users[username]
	if !ok {
		return errors.New("Could not find twith user: " + err.Error())
	}
	_, err = t.Client.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: broadcaster_id,
		//GameID:              "456789",
		//BroadcasterLanguage: "en",
		Title: title,
		Delay: 0,
	})
	if err != nil {
		return errors.New("Error changing twitch stream title: " + err.Error())
	}
	return nil
}

func (t *Twitch) GetUser(usernames []string) (map[string]string, error) {
	usersResp, err := t.Client.GetUsers(&helix.UsersParams{
		Logins: usernames,
	})
	if err != nil {
		panic(err)
	}
	users := map[string]string{}
	for i := range usersResp.Data.Users {
		user := usersResp.Data.Users[i]
		users[user.DisplayName] = user.ID
	}
	return users, nil
}
