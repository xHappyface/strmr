package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.Unmarshal(reqBody, &data)
		media_record, err := h.database.GetMediaRecordingByID(data.RecordingID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if media_record == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		yt_data, err := h.convertToYouTubeMetadata(*media_record)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		description := youtube.CreateMetadataText(yt_data.Tasks, "Starting stream")
		subtitle_file := youtube.CreateSubtitleText(yt_data.Subtitles)
		title := youtube.CreateTitleText(yt_data.Titles, "/")
		fmt.Println(description)
		fmt.Println(subtitle_file)
		fmt.Println(title)
		partial_file := media_record.Directory + "/" + media_record.FileName
		subtitle_file_name := partial_file + ".srt"
		err = ioutil.WriteFile(subtitle_file_name, []byte(subtitle_file), 0666)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		video_id, err := h.youtube.UploadVideo(partial_file+"."+media_record.Extension, title, description)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if video_id != nil {
			err = h.youtube.InsertCaption(*video_id, subtitle_file_name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
