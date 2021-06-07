package video

import (
	"net/http"

	"github.com/suconghou/videoproxy/request"
	"github.com/suconghou/youtubevideoparser"
)

// info 解析存在缓存,此处ProxyCall也缓存
func outPutTimedText(w http.ResponseWriter, r *http.Request, info *youtubevideoparser.VideoInfo) error {
	var (
		lang = r.URL.Query().Get("lang")
		url  = ""
	)
	for _, item := range info.Captions {
		url = item.URL
		if lang == "" || item.LanguageCode == lang {
			break
		}
	}
	if url == "" {
		http.Error(w, "lang not found", http.StatusNotFound)
		return nil
	}
	return request.ProxyCall(w, url, videoClient, r.Header)
}
