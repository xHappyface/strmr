package twitch

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

type Twitch struct {
	Client       *helix.Client
	Code         string
	Token        string
	RefreshToken string
}

type Category struct {
	ID        string
	BoxArtUrl string
}

func New(client *helix.Client) *Twitch {
	return &Twitch{
		Client:       client,
		Code:         "",
		Token:        "",
		RefreshToken: "",
	}
}

func (t *Twitch) ChangeStream(username string, title string, category_id string, tags []string) error {
	authorized, _, err := t.Client.ValidateToken(t.Token)
	if !authorized {
		return errors.New("Not authorized to change stream title: " + err.Error())
	}
	users, err := t.GetUser([]string{username})
	broadcaster_id, ok := users[username]
	if !ok {
		return errors.New("Could not find twith user: " + err.Error())
	}
	_, err = t.Client.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID:       broadcaster_id,
		GameID:              category_id,
		BroadcasterLanguage: "en",
		Title:               title,
		Tags:                tags,
		Delay:               0,
	})
	if err != nil {
		return errors.New("Error changing twitch stream title: " + err.Error())
	}
	return nil
}

func (t *Twitch) GetUser(usernames []string) (map[string]string, error) {
	authorized, _, err := t.Client.ValidateToken(t.Token)
	if !authorized {
		return nil, errors.New("Not authorized to get user: " + err.Error())
	}
	usersResp, err := t.Client.GetUsers(&helix.UsersParams{
		Logins: usernames,
	})
	if err != nil {
		return nil, errors.New("Could not get users: " + err.Error())
	}
	users := map[string]string{}
	for i := range usersResp.Data.Users {
		user := usersResp.Data.Users[i]
		users[user.DisplayName] = user.ID
	}
	return users, nil
}

func (t *Twitch) GetGames(names []string) (map[string]string, error) {
	authorized, _, err := t.Client.ValidateToken(t.Token)
	if !authorized {
		return nil, errors.New("Not authorized to get games")
	}
	if err != nil {
		return nil, errors.New("Error getting games in GetGames: " + err.Error())
	}
	gamesResp, err := t.Client.GetGames(&helix.GamesParams{
		Names: names,
	})
	if err != nil {
		return nil, errors.New("Could not get users: " + err.Error())
	}
	games := map[string]string{}
	for i := range gamesResp.Data.Games {
		game := gamesResp.Data.Games[i]
		games[game.Name] = game.ID
	}
	return games, nil
}

func (t *Twitch) GetUsers(names []string) (map[string]string, error) {
	authorized, _, err := t.Client.ValidateToken(t.Token)
	if !authorized {
		return nil, errors.New("Not authorized to get user")
	}
	if err != nil {
		return nil, errors.New("Error getting user in GetUsers: " + err.Error())
	}
	usersResp, err := t.Client.GetUsers(&helix.UsersParams{
		Logins: names,
	})
	if err != nil {
		return nil, errors.New("Could not get users: " + err.Error())
	}
	users := map[string]string{}
	for i := range usersResp.Data.Users {
		user := usersResp.Data.Users[i]
		users[user.Login] = user.ID
	}
	return users, nil
}

type Channel struct {
	Title        string
	CategoryName string
	CategoryID   string
	Tags         []string
}

func (t *Twitch) GetChannelInformation(ids []string) (map[string]Channel, error) {
	authorized, _, err := t.Client.ValidateToken(t.Token)
	if !authorized {
		return nil, errors.New("Not authorized to get channels")
	}
	if err != nil {
		return nil, errors.New("Error getting channels in SearchChannels: " + err.Error())
	}
	channelsResp, err := t.Client.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: ids,
	})
	if err != nil {
		return nil, errors.New("Could not get channels: " + err.Error())
	}
	channels := map[string]Channel{}
	for k := range channelsResp.Data.Channels {
		channel := Channel{
			Title:        channelsResp.Data.Channels[k].Title,
			CategoryName: channelsResp.Data.Channels[k].GameName,
			CategoryID:   channelsResp.Data.Channels[k].GameID,
			Tags:         channelsResp.Data.Channels[k].Tags,
		}
		channels[channelsResp.Data.Channels[k].BroadcasterName] = channel
	}
	return channels, nil
}

func (t *Twitch) SearchCategories(query string) (map[string]Category, error) {
	authorized, _, err := t.Client.ValidateToken(t.Token)
	if !authorized {
		return nil, errors.New("Not authorized to get categories")
	}
	if err != nil {
		return nil, errors.New("Error getting categories in SearchCategories: " + err.Error())
	}
	categoriesResp, err := t.Client.SearchCategories(&helix.SearchCategoriesParams{
		Query: query,
	})
	if err != nil {
		return nil, errors.New("Could not search categories: " + err.Error())
	}
	categories := map[string]Category{}
	for i := range categoriesResp.Data.Categories {
		category := categoriesResp.Data.Categories[i]
		categories[category.Name] = Category{
			ID:        category.ID,
			BoxArtUrl: category.BoxArtURL,
		}
	}
	return categories, nil
}
