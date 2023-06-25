package handlers

import (
	"errors"
	"log"
	"math"
	"net/http"
	"text/template"

	"github.com/jnrprgmr/strmr/pkg/database"
	"github.com/jnrprgmr/strmr/pkg/youtube"
)

type WrappedMediaRecording struct {
	database.MediaRecording
	Metadata youtube.YouTubeData
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
			Start: s.InsertTime - recording.StartTime - int64(math.Ceil(s.Duration)),
			End:   s.InsertTime - recording.StartTime,
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
		if *recording.EndTime < m.InsertTime {
			return nil, errors.New("recording end time is outside the medatada time")
		}
		yt_m := youtube.Metadata{
			Text:  m.MetadataValue,
			Start: m.InsertTime - recording.StartTime,
		}
		yt_metadata = append(yt_metadata, yt_m)
	}
	return yt_metadata, nil
}

func (h *Handlers) convertToYouTubeMetadata(media_recordings database.MediaRecording) (*youtube.YouTubeData, error) {
	tasks, err := h.database.GetMetadataByKeyAndTimeRange("task", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	yt_tasks, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, tasks)
	if err != nil {
		return nil, err
	}
	categories, err := h.database.GetMetadataByKeyAndTimeRange("category", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	initial_category, err := h.database.GetLatestMetadataByKeyBeforeTime("category", media_recordings.StartTime, 1)
	if err != nil {
		return nil, err
	}
	if len(initial_category) != 0 {

		categories = append(categories, initial_category...)
	}
	yt_categories, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, categories)
	if err != nil {
		return nil, err
	}
	tags, err := h.database.GetMetadataByKeyAndTimeRange("tags", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	initial_tags, err := h.database.GetLatestMetadataByKeyBeforeTime("tags", media_recordings.StartTime, 1)
	if err != nil {
		return nil, err
	}
	if len(initial_tags) != 0 {
		tags = append(tags, initial_tags...)
	}
	yt_tags, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, tags)
	if err != nil {
		return nil, err
	}

	descriptions, err := h.database.GetMetadataByKeyAndTimeRange("description", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	initial_descriptions, err := h.database.GetLatestMetadataByKeyBeforeTime("description", media_recordings.StartTime, 1)
	if err != nil {
		return nil, err
	}
	if len(initial_descriptions) != 0 {
		descriptions = append(descriptions, initial_descriptions...)
	}
	yt_descriptions, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, descriptions)
	if err != nil {
		return nil, err
	}

	titles, err := h.database.GetMetadataByKeyAndTimeRange("title", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	initial_title, err := h.database.GetLatestMetadataByKeyBeforeTime("title", media_recordings.StartTime, 1)
	if err != nil {
		return nil, err
	}
	if len(initial_title) != 0 {
		titles = append(titles, initial_title...)
	}
	yt_titles, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, titles)
	if err != nil {
		return nil, err
	}
	subtitles, err := h.database.GetSubtitlesByTimeRange(media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	yt_subtitles, err := ConvertRecordingSubtitlesToYouTubeSubtitles(media_recordings, subtitles)
	if err != nil {
		return nil, err
	}
	yt_data := &youtube.YouTubeData{
		Categories:   yt_categories,
		Titles:       yt_titles,
		Tasks:        yt_tasks,
		Subtitles:    yt_subtitles,
		Descriptions: yt_descriptions,
		Tags:         yt_tags,
	}
	return yt_data, nil
}

func (h *Handlers) YouTubeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		media_recordings, err := h.database.GetAllMediaRecordingsByUploaded(false) // get all non-uploaded videos
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cats, err := h.database.GetAllCategories()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		yt_cats, err := h.youtube.GetCategories()
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		yt_playlists, err := h.youtube.GetPlaylists()
		if err != nil {
			log.Fatalf("Cannot get playlists:%+v", err)
		}
		yt_playlists = append([]youtube.Playlist{
			{
				ID:   "",
				Name: "None",
			},
		}, yt_playlists...)
		data := map[int64]youtube.YouTubeData{}
		for i := range media_recordings {
			yt_data, err := h.convertToYouTubeMetadata(media_recordings[i])
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data[media_recordings[i].ID] = *yt_data
		}
		wrapped_recordings := []WrappedMediaRecording{}
		for i := range media_recordings {
			media_recording := media_recordings[i]
			wr := WrappedMediaRecording{media_recording, data[media_recording.ID]}
			wrapped_recordings = append(wrapped_recordings, wr)
		}
		tmpl := template.Must(template.ParseFiles("./templates/youtube.html"))
		tmpl.Execute(w, struct {
			Title             string
			Javascript        []string
			CSS               []string
			Recordings        []WrappedMediaRecording
			YouTubeCategories []youtube.Category
			YouTubePlaylists  []youtube.Playlist
			Categories        []database.Category
		}{
			Title: "YouTube settings",
			Javascript: []string{
				"vendor/jquery/jquery-3.6.3.min",
				"vendor/popper/popper-1.12.9.min",
				"vendor/bootstrap/bootstrap-4.0.0.min",
				"youtube",
			},
			CSS: []string{
				"vendor/bootstrap/bootstrap-4.0.0.min",
				"youtube",
			},
			Recordings:        wrapped_recordings,
			YouTubeCategories: yt_cats,
			YouTubePlaylists:  yt_playlists,
			Categories:        cats,
		})
	}
}
