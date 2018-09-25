package middleware

import (
	"dropboxshare/streampipe"
	"dropboxshare/util"
	"fmt"
	"net/http"

	"github.com/suconghou/youtubeVideoParser"
)

// ServeYoutubeImage proxy img
func ServeYoutubeImage(w http.ResponseWriter, r *http.Request, url string) error {
	return streampipe.Pipe(w, r, url, nil)
}

// GetYoutubeImageURL return img url
func GetYoutubeImageURL(quality string, id string, ext string) string {
	return youtubeVideoParser.GetYoutubeImageURL(id, ext, quality)
}

// ServeYoutubeVideo do video proxy
func ServeYoutubeVideo(w http.ResponseWriter, r *http.Request, url string) error {
	return streampipe.Pipe(w, r, url, nil)
}

// GetYoutubeVideoURL do get video url by vid
func GetYoutubeVideoURL(quality string, id string, ext string) (string, *youtubeVideoParser.VideoInfo, error) {
	info, err := youtubeVideoParser.Parse(id)
	util.Logger.Print(info)
	if err != nil {
		return "", nil, err
	}
	if ext == "json" {
		return "", info, nil
	}
	// url, _, err := info.GetStream(quality, ext)
	return "", nil, fmt.Errorf("no work")
}
