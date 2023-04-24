package youtube

import (
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	service *youtube.Service
}

func New(svc *youtube.Service) *YouTube {
	return &YouTube{
		service: svc,
	}
}

func (yt *YouTube) UploadVideo(file_name string) error {
	fmt.Println("starting video upload")
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:                "Test Vid",
			Description:          "0:00 Intro\n0:10 End\nRecorded: Now()",
			CategoryId:           "22",
			DefaultAudioLanguage: "en",
			DefaultLanguage:      "en",
			Tags: []string{
				"test",
				"vid",
			},
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "private"},
	}
	call := yt.service.Videos.Insert([]string{"snippet", "status"}, upload)
	file, err := os.Open(file_name)
	if err != nil {
		return errors.New("Error opening " + file_name + ": " + err.Error())
	}
	defer file.Close()
	response, err := call.Media(file).Do()
	if err != nil {
		return errors.New(err.Error())
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
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
