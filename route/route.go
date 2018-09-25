package route

import (
	"dropboxshare/middleware"
	"dropboxshare/util"
	"fmt"
	"net/http"
	"regexp"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

// RoutePath for all route
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/(files|imgs)/[\w\-\/\.]{0,100}$`), files},
	{regexp.MustCompile(`^/(large|medium|small)/([\w\-]{6,12})\.(mp4|flv|webm|3gp|json)$`), youtubeVideo},
	{regexp.MustCompile(`^/(large|medium|small)/([\w\-]{6,12})\.(jpg|webp)$`), youtubeImage},
}

func files(w http.ResponseWriter, r *http.Request, match []string) {
	err := middleware.Serve(w, r, match[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
	}
}

func youtubeVideo(w http.ResponseWriter, r *http.Request, match []string) {
	url, data, err := middleware.GetYoutubeVideoURL(match[1], match[2], match[3])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return
	}
	if url != "" {
		err = middleware.ServeYoutubeVideo(w, r, url)
		if err != nil {
			util.Logger.Print(err)
		}
		return
	}
	util.JSONPut(w, data)
}

func youtubeImage(w http.ResponseWriter, r *http.Request, match []string) {
	var url string = middleware.GetYoutubeImageURL(match[1], match[2], match[3])
	err := middleware.ServeYoutubeImage(w, r, url)
	if err != nil {
		util.Logger.Print(err)
	}
}
