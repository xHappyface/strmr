package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/jnrprgmr/strmr/internal/rest/handlers"
	"github.com/jnrprgmr/strmr/pkg/database"
	"github.com/jnrprgmr/strmr/pkg/obs"
	"github.com/jnrprgmr/strmr/pkg/twitch"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicklaw5/helix/v2"

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

func main() {
	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
