package handlers

import (
	"errors"
	"math"
	"net/http"
	"text/template"

	"github.com/jnrprgmr/strmr/pkg/database"
	"github.com/jnrprgmr/strmr/pkg/youtube"
)

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
	initial_category, err := h.database.GetLatestMetadataByKeyBeforeTime("category", media_recordings.StartTime)
	if err != nil {
		return nil, err
	}
	if initial_category != nil {
		categories = append(categories, *initial_category)
	}
	yt_categories, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, categories)
	if err != nil {
		return nil, err
	}
	tags, err := h.database.GetMetadataByKeyAndTimeRange("tags", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	initial_tags, err := h.database.GetLatestMetadataByKeyBeforeTime("tags", media_recordings.StartTime)
	if err != nil {
		return nil, err
	}
	if initial_tags != nil {
		tags = append(tags, *initial_tags)
	}
	yt_tags, err := ConvertRecordingMetadataToYouTubeMetadata(media_recordings, tags)
	if err != nil {
		return nil, err
	}
	titles, err := h.database.GetMetadataByKeyAndTimeRange("title", media_recordings.StartTime, *media_recordings.EndTime)
	if err != nil {
		return nil, err
	}
	initial_title, err := h.database.GetLatestMetadataByKeyBeforeTime("title", media_recordings.StartTime)
	if err != nil {
		return nil, err
	}
	if initial_title != nil {
		titles = append(titles, *initial_title)
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
		Categories: yt_categories,
		Titles:     yt_titles,
		Tasks:      yt_tasks,
		Subtitles:  yt_subtitles,
		Tags:       yt_tags,
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
		data := map[int64]youtube.YouTubeData{}
		for i := range media_recordings {
			yt_data, err := h.convertToYouTubeMetadata(media_recordings[i])
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data[media_recordings[i].ID] = *yt_data
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
				"youtube",
			},
			CSS: []string{
				"youtube",
			},
			Recordings: media_recordings,
		})
	}
}
