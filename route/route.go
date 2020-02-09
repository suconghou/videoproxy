package route

import (
	"net/http"
	"regexp"
	"videoproxy/video"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

// RoutePath for all route
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/video/([\w\-]{6,12})\.json$`), video.GetInfo},
	{regexp.MustCompile(`^/video/([\w\-]{6,12})/(\d{1,3})\.(mp4|webm)$`), video.ProxyOne},
	{regexp.MustCompile(`^/video/([\w\-]{6,12})/(\d{1,3})/(\d+-\d+)\.ts$`), video.ProxyPart},
	{regexp.MustCompile(`^/video/([\w\-]{6,12})\.(jpg|webp)$`), video.Image},

	{regexp.MustCompile(`^/video/api/(v3/videos)$`), video.Videos},
	{regexp.MustCompile(`^/video/api/(v3/search)$`), video.Search},
	{regexp.MustCompile(`^/video/api/(v3/channels)$`), video.Channels},
	{regexp.MustCompile(`^/video/api/(v3/playlists)$`), video.Playlists},
	{regexp.MustCompile(`^/video/api/(v3/playlistItems)$`), video.PlaylistItems},
	{regexp.MustCompile(`^/video/api/(v3/videoCategories)$`), video.Categories},
}
