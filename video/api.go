package video

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"videoproxy/request"
)

var (
	key = os.Getenv("YOUTUBE_API_KEY")
)

const (
	baseURL = "https://www.googleapis.com/youtube/%s"
)

// Videos proxy api to get video info , ?id=...
func Videos(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id, snippet, contentDetails, player, statistics, status")
	return call(w, r, match[1], q)
}

// Search proxy api , ?q=keyword&type=video&order=..&channelId=..
func Search(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id, snippet")
	q.Set("maxResults", "10")
	return call(w, r, match[1], q)
}

// Channels proxy api , ?forUsername=.. / ?id=..
func Channels(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id,snippet,contentDetails,statistics,invideoPromotion")
	return call(w, r, match[1], q)
}

// Playlists proxy api , ?id=.. / ?channelId=..
func Playlists(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id, snippet, status")
	return call(w, r, match[1], q)
}

// PlaylistItems proxy api , ?playlistId=..
func PlaylistItems(w http.ResponseWriter, r *http.Request, match []string) error {
	var q = r.URL.Query()
	q.Set("part", "id, snippet, contentDetails, status")
	q.Set("maxResults", "50")
	return call(w, r, match[1], q)
}

func call(w http.ResponseWriter, r *http.Request, t string, q url.Values) error {
	q.Set("key", key)
	var url = fmt.Sprintf(baseURL, t) + "?" + q.Encode()
	return request.ProxyCall(w, url)
}
