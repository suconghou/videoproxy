package video

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/suconghou/videoproxy/request"
	"github.com/suconghou/videoproxy/util"

	"github.com/suconghou/youtubevideoparser"
)

var (
	preferList          = []string{"18", "59", "22", "37", "243", "134", "244", "135", "247", "136", "248", "137", "242", "133", "278", "160"}
	imageClient         = util.MakeClient("IMAGE_PROXY", time.Minute)
	videoClient         = util.MakeClient("VIDEO_PROXY", time.Minute)
	youtubeImageHostMap = map[string]string{
		"jpg":  "http://i.ytimg.com/vi/",
		"webp": "http://i.ytimg.com/vi_webp/",
	}
	r1 int
	r2 int
)

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func init() {
	codepass := os.Getenv("CODE_PASS")
	if strings.Contains(codepass, ",") {
		arr := strings.Split(codepass, ",")
		var err error
		if r1, err = strconv.Atoi(arr[0]); err == nil {
			r2, _ = strconv.Atoi(arr[1])
		}
	}
}

func getinfo(id string) (*youtubevideoparser.VideoInfo, error) {
	parser, err := youtubevideoparser.NewParser(id, videoClient)
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
	return request.Pipe(w, r, url, imageClient, nil)
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
	_, err = util.JSONPut(w, info, http.StatusOK, 864000)
	return err
}

// ProxyAuto find playable a&v stream
func ProxyAuto(w http.ResponseWriter, r *http.Request, match []string) error {
	var (
		query     = r.URL.Query()
		info, err = getinfo(match[1])
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	var s = findItem(info, query.Get("prefer"))
	if s == nil {
		http.NotFound(w, r)
		return nil
	}
	var filename = ""
	if query.Get("download") == "1" {
		if strings.Contains(s.Type, "mp4") {
			filename = fmt.Sprintf("%s.%s", info.Title, "mp4")
		} else {
			filename = fmt.Sprintf("%s.%s", info.Title, "webm")
		}
	}
	return request.Pipe(w, r, s.URL, videoClient, func(res, to http.Header) {
		if filename != "" {
			name := url.PathEscape(filename)
			to.Set("Content-Disposition", fmt.Sprintf("attachment;filename* = UTF-8''%s", name))
		}
	})
}

func findItem(info *youtubevideoparser.VideoInfo, prefers string) *youtubevideoparser.StreamItem {
	for _, itag := range strings.Split(prefers, ",") {
		if v, ok := info.Streams[itag]; ok {
			if v.ContentLength != "" {
				return v
			}
		}
	}
	for _, itag := range preferList {
		if v, ok := info.Streams[itag]; ok {
			if v.ContentLength != "" {
				return v
			}
		}
	}
	for _, v := range info.Streams {
		if v.ContentLength != "" {
			return v
		}
	}
	return nil
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
		return request.Pipe(w, r, s.URL, videoClient, nil)
	}
	return request.ProxyData(w, r, s.URL+"&range="+ts, videoClient)
}

// AuthCode decode vid if encoded
func AuthCode(handler func(http.ResponseWriter, *http.Request, []string) error) func(http.ResponseWriter, *http.Request, []string) error {
	if r1 > 0 && r2 > 0 {
		return func(w http.ResponseWriter, r *http.Request, match []string) error {
			vid, err := util.DecodeVid(match[1], r1, r2)
			if err != nil {
				http.Error(w, "bad request", http.StatusForbidden)
				return err
			}
			match[1] = vid
			return handler(w, r, match)
		}
	}
	return handler
}
