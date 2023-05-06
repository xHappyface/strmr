package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jnrprgmr/strmr/pkg/youtube"
)

type YouTubeUpload struct {
	RecordingID int64 `json:"recording_id"`
}

func (h *Handlers) YouTubeUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data YouTubeUpload
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		media_record, err := h.database.GetMediaRecordingByID(data.RecordingID)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if media_record == nil {
			h.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		yt_data, err := h.convertToYouTubeMetadata(*media_record)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		description := ""
		categories := []string{}
		for i := range yt_data.Categories {
			categories = append(categories, yt_data.Categories[i].Text)
		}
		tags := []string{}
		unique_tags := map[string]bool{}
		for i := range yt_data.Tags {
			real_tags := strings.Split(yt_data.Tags[i].Text, ",")
			for j := range real_tags {
				unique_tags[real_tags[j]] = true
			}
		}
		for tag := range unique_tags {
			tags = append(tags, "#"+tag)
		}
		recording_time := time.Unix(media_record.StartTime, 0).UTC().Format(time.RFC3339Nano)
		description = description + "Categories:\n" + strings.Join(categories, "\n") + "\n\n"
		description = description + "Timestamps:\n" + youtube.CreateMetadataText(yt_data.Tasks, "Starting stream") + "\n"
		description = description + strings.Join(tags, " ") + "\n"
		description = description + "Streamed: " + recording_time
		subtitle_file := youtube.CreateSubtitleText(yt_data.Subtitles)
		title := youtube.CreateTitleText(yt_data.Titles, "/")
		partial_file := media_record.Directory + "/" + media_record.FileName
		subtitle_file_name := partial_file + ".srt"
		err = ioutil.WriteFile(subtitle_file_name, []byte(subtitle_file), 0666)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		video_id, err := h.youtube.UploadVideo(partial_file+"."+media_record.Extension, title, description, tags, recording_time)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if video_id != nil {
			err = h.youtube.InsertCaption(*video_id, subtitle_file_name)
			if err != nil {
				h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		err = h.database.SetMediaRecordingUploadedByID(data.RecordingID, true)
		if err != nil {
			h.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
