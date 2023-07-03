package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/jnrprgmr/strmr/internal/rest/handlers"
	"github.com/jnrprgmr/strmr/pkg/brdcstr"
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
	Brdcstr  brdcstr.Config  `yaml:"brdcstr"`
}

func loadConfig() (*Config, error) {
	c := Config{}
	yamlFile, err := os.ReadFile("conf/local.yaml")
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v", err)
	}
	if err = yaml.Unmarshal(yamlFile, &c); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config file: %v", err)
	}
	return &c, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Could not open OAUTH file, please run the scripts/auth.py to get the authorization json: %v", err)
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
	obs := obs.New(obsCli, "strmr-screen", "strmr-task-text", "strmr-task-background", "strmr-avatar", "strmr-overlay-text", "strmr-overlay-background")
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
		log.Fatalf("Error making twitch client: %v", err)
	}
	brdcstr_client := brdcstr.New(c.Brdcstr.Host, c.Brdcstr.Port)
	_, err = brdcstr_client.Alive()
	if err != nil {
		log.Fatalf("Error checking brdcstr alive: %v", err)
	}
	ctx := context.Background()
	// Used this video to help setup google coud project and get client secrets https://www.youtube.com/watch?v=aFwZgth790Q
	// run auth.py to generate oauth token and allow youtube API
	b, err := os.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, youtubeApi.YoutubeUploadScope, youtubeApi.YoutubeForceSslScope, youtubeApi.YoutubepartnerScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	f, err := os.Open(os.Getenv("GOOGLE_OAUTH_TOKENS"))
	if err != nil {
		log.Fatalf("Could not open OAUTH file, please run the scripts/auth.py to get the authorization json: %v", err)
	}
	defer f.Close()
	t := &oauth2.Token{}
	if err = json.NewDecoder(f).Decode(t); err != nil {
		log.Fatal(err)
	}
	tok, err := t, err
	if err != nil {
		log.Fatal(err)
	}
	client := config.Client(ctx, tok)
	service, err := youtubeApi.New(client)
	if err != nil {
		log.Fatal(err)
	}
	yt := youtube.New(service)
	twitch := twitch.New(twitchCli)
	h := handlers.New(twitch, obs, yt, db)
	http.HandleFunc("/twitch", h.TwitchHandler)
	http.HandleFunc("/twitch/update", h.TwitchUpdateHandler)
	http.HandleFunc("/twitch/auth", h.TwitchAuthHandler)
	http.HandleFunc("/twitch/search/categories", h.TwitchSearchCategoriesHandler)

	http.HandleFunc("/obs", h.ObsHandler)
	http.HandleFunc("/obs/task", h.UpdateOBSTask)
	http.HandleFunc("/obs/scene/create", h.CreateScene)
	http.HandleFunc("/obs/stream", h.UpdateOBSStream)
	http.HandleFunc("/obs/overlay", h.UpdateOBSOverlay)

	http.HandleFunc("/youtube", h.YouTubeHandler)
	http.HandleFunc("/youtube_upload", h.YouTubeUploadHandler)
	http.HandleFunc("/youtube_category", h.YouTubeCategoryHandler)

	http.HandleFunc("/avatar_status", h.AvatarStatus)
	http.HandleFunc("/avatar", h.Avatar)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	background_config_metadata, err := db.GetLatestMetadataByKey("task_background_config", 1)
	if err != nil || len(background_config_metadata) != 1 {
		log.Fatalf("Cannot get background config metadata:%+v", err)
	}
	task_config_metadata, err := db.GetLatestMetadataByKey("task_config", 1)
	if err != nil || len(task_config_metadata) != 1 {
		log.Fatalf("Cannot get task config metadata:%+v", err)
	}
	task_metadata, err := db.GetLatestMetadataByKey("task", 1)
	if err != nil || len(task_metadata) != 1 {
		log.Fatalf("Cannot get task metadata:%+v", err)
	}
	err = obs.RefreshSources(background_config_metadata[0], task_metadata[0].MetadataValue, task_config_metadata[0])
	if err != nil {
		fmt.Println(err)
	}
	s := &http.Server{
		Addr: "0.0.0.0:8080",
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen: %v\n", err)
		}
	}()
	fmt.Println("Server Started")

	<-done
	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server Exited Properly")
}
