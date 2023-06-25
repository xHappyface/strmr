package youtube

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	service *youtube.Service
}

type Playlist struct {
	ID   string
	Name string
}

type Metadata struct {
	Text  string
	Start int64
}

func GetTimestamp(ts int64) string {
	hours := ts / 3600
	rem := ts % 3600
	min := rem / 60
	rem = rem % 60
	sec := rem
	return fmt.Sprintf("%02s", strconv.Itoa(int(hours))) + ":" + fmt.Sprintf("%02s", strconv.Itoa(int(min))) + ":" + fmt.Sprintf("%02s", strconv.Itoa(int(sec)))
}

type Subtitle struct {
	Text  string
	Start int64
	End   int64
}

func CreateSubtitleText(subtitles []Subtitle) string {
	text := ""
	for i := range subtitles {
		s := subtitles[i]
		text = text + strconv.Itoa(i+1) + "\n" + GetTimestamp(s.Start) + ",000 --> " + GetTimestamp(s.End) + ",000\n" + s.Text + "\n\n"
	}
	return text
}

type YouTubeData struct {
	File         string
	Categories   []Metadata
	Titles       []Metadata
	Tags         []Metadata
	Descriptions []Metadata
	Tasks        []Metadata
	Subtitles    []Subtitle
}

func CreateMetadataText(metadata []Metadata, initial string) string {
	text := "00:00:00 " + initial + "\n"
	for i := range metadata {
		text = text + GetTimestamp(metadata[i].Start) + " " + metadata[i].Text + "\n"
	}
	return text
}

func CreateTitleText(titles []Metadata, delimiter string) string {
	text := []string{}
	for i := range titles {
		text = append(text, titles[i].Text)
	}
	return strings.Join(text, delimiter)
}

func New(svc *youtube.Service) *YouTube {
	return &YouTube{
		service: svc,
	}
}

func (yt *YouTube) UploadVideo(file_name string, title string, description string, tags []string, recording_time string, category string) (*string, error) {
	fmt.Println("starting video upload")
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:                title,
			Description:          description,
			CategoryId:           category,
			DefaultAudioLanguage: "en",
			DefaultLanguage:      "en",
			Tags:                 tags,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus:           "private",
			SelfDeclaredMadeForKids: false,
			ForceSendFields: []string{
				"SelfDeclaredMadeForKids",
			},
		},
		RecordingDetails: &youtube.VideoRecordingDetails{
			RecordingDate: recording_time,
		},
	}
	call := yt.service.Videos.Insert([]string{"snippet", "status", "recordingDetails"}, upload)
	file, err := os.Open(file_name)
	if err != nil {
		return nil, errors.New("Error opening " + file_name + ": " + err.Error())
	}
	defer file.Close()
	response, err := call.Media(file).Do()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
	return &response.Id, nil
}

func (yt *YouTube) GetPlaylists() ([]Playlist, error) {
	call := yt.service.Playlists.List([]string{
		"snippet",
		"id",
	})
	resp, err := call.Mine(true).Do()
	if err != nil {
		return nil, errors.New("Cannot get youtube playlists: " + err.Error())
	}
	playlists := []Playlist{}
	for i := range resp.Items {
		playlists = append(playlists, Playlist{
			ID:   resp.Items[i].Id,
			Name: resp.Items[i].Snippet.Title,
		})
	}
	return playlists, nil
}

func (yt *YouTube) InsertPlaylist(video_id string, playlist_id string) error {
	fmt.Println("Adding to playlist")
	playlist := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlist_id,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: video_id,
			},
		},
	}
	response, err := yt.service.PlaylistItems.Insert([]string{"snippet"}, playlist).Do()
	if err != nil {
		return errors.New("Could not add video to playlist: " + err.Error())
	}
	fmt.Printf("added vide to playlist successful! ID: %v\n", response.Id)
	return nil
}

func (yt *YouTube) InsertCaption(video_id string, file_name string) error {
	fmt.Println("starting caption upload")
	caption := &youtube.Caption{
		Snippet: &youtube.CaptionSnippet{
			Language: "en",
			Name:     "English subtitles",
			VideoId:  video_id,
		},
	}
	call := yt.service.Captions.Insert([]string{"snippet"}, caption)
	file, err := os.Open(file_name)
	if err != nil {
		return errors.New("Error opening " + file_name + ": " + err.Error())
	}
	defer file.Close()
	response, err := call.Media(file).Do()
	if err != nil {
		return errors.New(err.Error())
	}
	fmt.Printf("Upload captions successful! ID: %v\n", response.Id)
	return nil
}

type Category struct {
	ID    string
	Title string
}

func (yt *YouTube) GetCategories() ([]Category, error) {
	resp, err := yt.service.VideoCategories.List([]string{"id"}).RegionCode("US").Do()
	if err != nil {
		return nil, errors.New("Cannot get youtube categories: " + err.Error())
	}
	categories := []Category{}
	for i := range resp.Items {
		categories = append(categories, Category{
			ID:    resp.Items[i].Id,
			Title: resp.Items[i].Snippet.Title,
		})
	}
	return categories, nil
}
