package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/jnrprgmr/strmr/pkg/database"
	"github.com/jnrprgmr/strmr/pkg/youtube"
)

type YouTubeData struct {
	Recording  database.MediaRecording
	Categories []database.Metadata
	Titles     []database.Metadata
	Tags       []database.Metadata
	Subtitles  []database.Subtitle
}

func ConvertRecordingSubtitlesToYouTubeSubtitles(recording database.MediaRecording, subtitles []database.Subtitle) ([]youtube.Subtitle, error) {
	if recording.EndTime == nil {
		return nil, errors.New("recording has not ended")
	}
	yt_subtitles := []youtube.Subtitle{}
	if len(subtitles) < 1 {
		return yt_subtitles, nil
	}
	for i := range subtitles {
		s := subtitles[i]
		if recording.StartTime > s.InsertTime || *recording.EndTime < s.InsertTime+int64(s.Duration) {
			return nil, errors.New("recording start or end time is outside the subtitle time")
		}
		yt_s := youtube.Subtitle{
			Text:  s.Subtitle,
			Start: s.InsertTime - recording.StartTime,
			End:   (s.InsertTime + int64(s.Duration)) - recording.StartTime,
		}
		yt_subtitles = append(yt_subtitles, yt_s)
	}
	return yt_subtitles, nil
}

func ConvertRecordingMetadataToYouTubeMetadata(recording database.MediaRecording, metadata []database.Metadata) ([]youtube.Metadata, error) {
	if recording.EndTime == nil {
		return nil, errors.New("recording has not ended")
	}
	yt_metadata := []youtube.Metadata{}
	if len(metadata) < 1 {
		return yt_metadata, nil
	}
	for i := range metadata {
		m := metadata[i]
		if recording.StartTime > m.InsertTime || *recording.EndTime < m.InsertTime {
			return nil, errors.New("recording start or end time is outside the medatada time")
		}
		yt_m := youtube.Metadata{
			Text:  m.MetadataValue,
			Start: m.InsertTime - recording.StartTime,
		}
		yt_metadata = append(yt_metadata, yt_m)
	}
	return yt_metadata, nil
}

func (h *Handlers) YouTubeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		media_recordings, err := h.database.GetAllMediaRecordingsByUploaded(false) // get all non-uploaded videos
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := map[int64]youtube.YouTubeData{}
		for i := range media_recordings {
			tasks, err := h.database.GetMetadataByKeyAndTimeRange("task", media_recordings[i].StartTime, *media_recordings[i].EndTime)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			yt_tasks, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings[i], tasks)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			categories, err := h.database.GetMetadataByKeyAndTimeRange("category", media_recordings[i].StartTime, *media_recordings[i].EndTime)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			yt_categories, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings[i], categories)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			titles, err := h.database.GetMetadataByKeyAndTimeRange("title", media_recordings[i].StartTime, *media_recordings[i].EndTime)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			yt_titles, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings[i], titles)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			subtitles, err := h.database.GetSubtitlesByTimeRange(media_recordings[i].StartTime, *media_recordings[i].EndTime)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			yt_subtitles, err := ConvertRecordingSubtitlesToYouTubeSubtitles(media_recordings[i], subtitles)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			yt_data := youtube.YouTubeData{
				Categories: yt_categories,
				Titles:     yt_titles,
				Tasks:      yt_tasks,
				Subtitles:  yt_subtitles,
			}
			fmt.Println(youtube.CreateMetadataText(yt_data.Categories, "Stream starting"))
			fmt.Println(youtube.CreateMetadataText(yt_data.Titles, "Stream starting"))
			fmt.Println(youtube.CreateMetadataText(yt_data.Tasks, "Stream starting"))
			fmt.Println(youtube.CreateSubtitleText(yt_data.Subtitles))
			data[media_recordings[i].ID] = yt_data
		}
		tmpl := template.Must(template.ParseFiles("./templates/youtube.html"))
		tmpl.Execute(w, struct {
			Title      string
			Javascript []string
			CSS        []string
			Recordings []database.MediaRecording
		}{
			Title: "OBS stream settings",
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"obs",
			},
			CSS: []string{
				"obs",
			},
			Recordings: media_recordings,
		})
	}
}
