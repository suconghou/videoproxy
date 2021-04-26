package video

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/suconghou/youtubevideoparser"
)

var (
	alist = "249,250,251,600,140,599"
	vlist = "247,136,244,135,243,134,242,133,278,160"
)

func formatDuration(t int) string {
	var text = "PT"
	var hour = t / 3600
	var min = (t % 3600) / 60
	var sec = (t % 60)
	if hour >= 1 {
		text += fmt.Sprintf("%dH", hour)
	}
	if min >= 1 {
		text += fmt.Sprintf("%dM", min)
	}
	text += fmt.Sprintf("%dS", sec)
	return text
}

func outPutMpd(w http.ResponseWriter, r *http.Request, info *youtubevideoparser.VideoInfo) error {
	xml, err := buildXML(r, info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	h := w.Header()
	h.Set("Content-Type", "application/dash+xml")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Max-Age", "864000")
	h.Set("Cache-Control", "public,max-age=864000")
	_, err = w.Write([]byte(xml))
	return err
}

func buildXML(r *http.Request, info *youtubevideoparser.VideoInfo) (string, error) {
	duration, err := strconv.Atoi(info.Duration)
	if err != nil {
		return "", err
	}
	var (
		query   = r.URL.Query()
		t       = formatDuration(duration)
		b       = strings.Builder{}
		header  = fmt.Sprintf("<MPD xmlns=\"urn:mpeg:dash:schema:mpd:2011\" profiles=\"urn:mpeg:dash:profile:isoff-on-demand:2011\" minBufferTime=\"PT2S\" mediaPresentationDuration=\"%s\" type=\"static\"><Period>", t)
		video   string
		audio   string
		patharr = strings.Split(strings.ReplaceAll(r.URL.Path, ".mpd", ""), "/")
	)
	var ID = patharr[len(patharr)-1]
	audio, video, err = buildItem(info, ID, duration, query.Get("a"), query.Get("v"))
	if err != nil {
		return "", err
	}
	b.WriteString(header)
	b.WriteString("<AdaptationSet segmentAlignment=\"true\">")
	b.WriteString(video)
	b.WriteString("</AdaptationSet>")
	b.WriteString("<AdaptationSet segmentAlignment=\"true\">")
	b.WriteString(audio)
	b.WriteString("</AdaptationSet>")
	b.WriteString("</Period></MPD>")
	return b.String(), nil
}

func buildItem(info *youtubevideoparser.VideoInfo, ID string, duration int, a string, v string) (string, string, error) {
	audio := findStream(info, a+","+alist, "audio")
	video := findStream(info, v+","+vlist, "video")
	if audio == nil || video == nil {
		return "", "", fmt.Errorf("failed to get video or audio")
	}
	astr, err := formatItem(ID, duration, audio)
	if err != nil {
		return "", "", err
	}
	vstr, err := formatItem(ID, duration, video)
	if err != nil {
		return "", "", err
	}
	return astr, vstr, nil
}

func formatItem(ID string, duration int, item *youtubevideoparser.StreamItem) (string, error) {
	len, err := strconv.Atoi(item.ContentLength)
	if err != nil {
		return "", err
	}
	var (
		ext        = "mp4"
		typeinfo   = strings.Split(item.Type, ";")
		mime       = typeinfo[0]
		codecs     = typeinfo[1]
		indexRange = fmt.Sprintf("%s-%s", item.IndexRange.Start, item.IndexRange.End)
		initRange  = fmt.Sprintf("%s-%s", item.InitRange.Start, item.InitRange.End)
		bandwidth  = 8 * (len / duration)
	)
	if strings.Contains(mime, "webm") {
		ext = "webm"
	}
	var (
		baseurl = fmt.Sprintf("%s/%s.%s", ID, item.Itag, ext)
		record  = fmt.Sprintf("<Representation id=\"%s\" bandwidth=\"%d\" %s mimeType=\"%s\"><BaseURL>%s</BaseURL><SegmentBase indexRange=\"%s\"><Initialization range=\"%s\"/></SegmentBase></Representation>", item.Itag, bandwidth, codecs, mime, baseurl, indexRange, initRange)
	)
	return record, nil
}

func findStream(info *youtubevideoparser.VideoInfo, prefers string, mime string) *youtubevideoparser.StreamItem {
	for _, itag := range strings.Split(prefers, ",") {
		if v, ok := info.Streams[itag]; ok {
			if v.ContentLength != "" && v.InitRange.Start != "" && v.IndexRange.Start != "" {
				return v
			}
		}
	}
	for _, v := range info.Streams {
		if v.ContentLength != "" && v.InitRange.Start != "" && v.IndexRange.Start != "" && strings.Contains(v.Type, mime) {
			return v
		}
	}
	return nil
}
