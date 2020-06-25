package video

import (
	"fmt"
	"net/http"

	"github.com/suconghou/videoproxy/request"
	"github.com/suconghou/videoproxy/util"

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

func getinfo(id string) (*youtubevideoparser.VideoInfo, error) {
	parser, err := youtubevideoparser.NewParser(id)
	if err != nil {
		return nil, err
	}
	return parser.Parse()
}

// Image proxy yputube image , default/mqdefault/hqdefault/sddefault/maxresdefault
func Image(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		id  = match[1]
		ext = match[2]
		url = fmt.Sprintf("%s%s/%s.%s", youtubeImageHostMap[ext], id, "mqdefault", ext)
	)
	return request.Pipe(w, r, url)
}

// GetInfo for info
func GetInfo(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		info, err = getinfo(match[1])
	)
	if err != nil {
		util.JSONPut(w, resp{-1, err.Error()}, http.StatusInternalServerError, 1)
		return err
	}
	// 为使接口长缓存,默认不出易失效数据
	if r.URL.Query().Get("info") != "all" {
		for _, s := range info.Streams {
			s.URL = ""
		}
	}
	_, err = util.JSONPut(w, info, http.StatusOK, 604800)
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
		info, err = getinfo(id)
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	s := info.Streams[itag]
	if s == nil {
		http.NotFound(w, r)
		return nil
	}
	if ts == "" {
		return request.Pipe(w, r, s.URL)
	}
	return request.ProxyData(w, r, s.URL+"&range="+ts)
}
