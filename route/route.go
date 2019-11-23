package route

import (
	"dropboxshare/video"
	"net/http"
	"regexp"
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

	{regexp.MustCompile(`^/image/([\w\-]{6,12})\.(jpg|webp)$`), video.Image},
}
