package route

import (
	"net/http"
	"regexp"

	"github.com/suconghou/videoproxy/video"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

// Route for all route
var Route = []routeInfo{
	{regexp.MustCompile(`^/video/([\w\-]{6,15})\.(json|mpd)$`), video.AuthCode(video.GetInfo)},
	{regexp.MustCompile(`^/video/([\w\-]{6,15})/(\d{1,3})\.(mp4|webm)$`), video.AuthCode(video.ProxyOne)},
	{regexp.MustCompile(`^/video/([\w\-]{6,15})/(\d{1,3})/(\d+-\d+)\.ts$`), video.AuthCode(video.ProxyPart)},
	{regexp.MustCompile(`^/video/([\w\-]{6,15})\.(jpg|webp)$`), video.AuthCode(video.Image)},
	{regexp.MustCompile(`^/video/([\w\-]{6,15})\.(mp4|webm)$`), video.AuthCode(video.ProxyAuto)},

	{regexp.MustCompile(`^/video/api/(v3/videos)$`), video.Videos},
	{regexp.MustCompile(`^/video/api/(v3/search)$`), video.Search},
	{regexp.MustCompile(`^/video/api/(v3/channels)$`), video.Channels},
	{regexp.MustCompile(`^/video/api/(v3/playlists)$`), video.Playlists},
	{regexp.MustCompile(`^/video/api/(v3/playlistItems)$`), video.PlaylistItems},
	{regexp.MustCompile(`^/video/api/(v3/videoCategories)$`), video.Categories},
}
