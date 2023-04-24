package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/jnrprgmr/strmr/internal/rest/handlers"
	"github.com/jnrprgmr/strmr/pkg/database"
	"github.com/jnrprgmr/strmr/pkg/obs"
	"github.com/jnrprgmr/strmr/pkg/twitch"
	"github.com/jnrprgmr/strmr/pkg/youtube"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicklaw5/helix/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	youtubeApi "google.golang.org/api/youtube/v3"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database database.Config `yaml:"db"`
	OBS      obs.Config      `yaml:"obs"`
}

func loadConfig() (*Config, error) {
	c := Config{}
	yamlFile, err := ioutil.ReadFile("conf/local.yaml")
	if err != nil {
		return nil, errors.New("Failed to read config file: " + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, errors.New("Failed to unmarshal config file: " + err.Error())
	}
	return &c, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New("could not open OAUTH file, please run the scripts/auth.py to get the authorization json: " + err.Error())
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func main() {
	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	obs_password := os.Getenv("OBS_PASSWORD")
	obsCli, err := goobs.New(c.OBS.Host+":"+c.OBS.Port, goobs.WithPassword(obs_password))
	if err != nil {
		log.Fatal(err)
	}
	defer obsCli.Disconnect()
	obs := obs.New(obsCli, "test", "background")
	sqlxConn, err := database.GetDB(c.Database.Name)
	if err != nil {
		log.Fatal(err)
	}
	db := database.New(sqlxConn)
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	twitchCli, err := helix.NewClient(&helix.Options{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURI:  "http://localhost:8080/twitch/auth",
	})
	if err != nil {
		panic("error making twitch client: " + err.Error())
	}
	ctx := context.Background()
	// Used this video to help setup google coud project and get client secrets https://www.youtube.com/watch?v=aFwZgth790Q
	// run auth.py to generate oauth token and allow youtube API
	b, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, youtubeApi.YoutubeUploadScope, youtubeApi.YoutubeForceSslScope, youtubeApi.YoutubepartnerScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	f, err := os.Open(os.Getenv("GOOGLE_OAUTH_TOKENS"))
	if err != nil {
		panic(errors.New("could not open OAUTH file, please run the scripts/auth.py to get the authorization json: " + err.Error()))
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	tok, err := t, err
	if err != nil {
		panic(err)
	}
	client := config.Client(ctx, tok)
	service, err := youtubeApi.New(client)
	if err != nil {
		log.Fatal(err)
	}
	// _ = youtube.New(service)
	// ////
	yt := youtube.New(service)
	err = yt.InsertCaption("8ICxJgs7zJk", "/media/jnrprgmr/7C000E4D000E0EB8/Videos/captions.srt")
	if err != nil {
		log.Fatal(err)
	}
	// /////
	twitch := twitch.New(twitchCli)
	h := handlers.New(twitch, obs, db)
	http.HandleFunc("/twitch", h.TwitchHandler)
	http.HandleFunc("/twitch/update", h.TwitchUpdateHandler)
	http.HandleFunc("/twitch/auth", h.TwitchAuthHandler)
	http.HandleFunc("/twitch/search/categories", h.TwitchSearchCategoriesHandler)

	http.HandleFunc("/obs", h.ObsHandler)
	http.HandleFunc("/obs/task", h.UpdateOBSTask)
	http.HandleFunc("/obs/scene/create", h.CreateScene)
	http.HandleFunc("/obs/stream", h.UpdateOBSStream)

	http.HandleFunc("/youtube", h.YouTubeHandler)

	http.HandleFunc("/avatar_status", h.AvatarStatus)
	http.HandleFunc("/avatar", h.Avatar)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
