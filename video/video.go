package video

import (
	"fmt"
	"net/http"
	"videoproxy/request"
	"videoproxy/util"

	"github.com/suconghou/youtubevideoparser"
)

var (
	youtubeImageHostMap = map[string]string{
		"jpg":  "http://i.ytimg.com/vi/",
		"webp": "http://i.ytimg.com/vi_webp/",
	}
)

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Image proxy yputube image , default/mqdefault/hqdefault/sddefault/maxresdefault
func Image(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		id  = match[1]
		ext = match[2]
		url = fmt.Sprintf("%s%s/%s.%s", youtubeImageHostMap[ext], id, "mqdefault", ext)
	)
	return request.Pipe(w, r, url, "")
}

// GetInfo for info
func GetInfo(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		info        *youtubevideoparser.VideoInfo
		id          = match[1]
		parser, err = youtubevideoparser.NewParser(id)
	)
	if err != nil {
		util.JSONPut(w, resp{-1, err.Error()})
		return err
	}
	if info, err = parser.Parse(); err != nil {
		util.JSONPut(w, resp{-2, err.Error()})
		return err
	}
	_, err = util.JSONPut(w, info)
	return err
}

// ProxyOne proxy whole video
func ProxyOne(w http.ResponseWriter, r *http.Request, match []string) error {
	return proxy(w, r, match[1], match[2], "")
}

// ProxyPart proxy a range part
func ProxyPart(w http.ResponseWriter, r *http.Request, match []string) error {
	return proxy(w, r, match[1], match[2], match[3])
}

// proxy proxy a range part
func proxy(w http.ResponseWriter, r *http.Request, id string, itag string, ts string) error {
	var (
		info        *youtubevideoparser.VideoInfo
		parser, err = youtubevideoparser.NewParser(id)
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if info, err = parser.Parse(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if v, has := info.Streams[itag]; has {
		return request.Pipe(w, r, v.URL, ts)
	}
	http.NotFound(w, r)
	return nil
}
