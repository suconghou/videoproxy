package video

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/suconghou/videoproxy/request"
	"github.com/suconghou/videoproxy/util"
)

var (
	key       = os.Getenv("YOUTUBE_API_KEY")
	apiClient = util.MakeClient("API_PROXY", time.Minute)
)

const (
	baseURL = "https://www.googleapis.com/youtube/%s"
)

// Videos proxy api to get video info , ?id=.. / ?chart=mostPopular&maxResults=20
func Videos(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet,contentDetails,statistics")
	return call(w, r, match[1], q)
}

// Search proxy api , ?q=keyword&type=video&order=..&channelId=..
func Search(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet")
	return call(w, r, match[1], q)
}

// Channels proxy api , ?forUsername=.. / ?id=.. / ?categoryId=..
func Channels(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet,contentDetails,statistics")
	return call(w, r, match[1], q)
}

// Playlists proxy api , ?id=.. / ?channelId=..
func Playlists(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet")
	return call(w, r, match[1], q)
}

// PlaylistItems proxy api , ?playlistId=..
func PlaylistItems(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet,contentDetails")
	return call(w, r, match[1], q)
}

// Categories proxy api , ?id=.. / ?regionCode=
func Categories(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet")
	return call(w, r, match[1], q)
}

func call(w http.ResponseWriter, r *http.Request, t string, q url.Values) error {
	if q.Get("key") == "" {
		q.Set("key", key)
	}
	var url = fmt.Sprintf(baseURL, t) + "?" + q.Encode()
	return request.ProxyCall(w, url, apiClient, r.Header, nil)
}
