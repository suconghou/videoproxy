package request

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	client     = &http.Client{Timeout: time.Minute}
	fwdHeaders = []string{
		"User-Agent",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"If-Modified-Since",
		"If-None-Match",
		"Range",
		"Content-Length",
		"Content-Type",
	}
	exposeHeaders = []string{
		"Accept-Ranges",
		"Content-Range",
		"Content-Length",
		"Content-Type",
		"Content-Encoding",
		"Date",
		"Expires",
		"Last-Modified",
		"Etag",
		"Cache-Control",
	}
)

// Pipe Proxy request
func Pipe(w http.ResponseWriter, r *http.Request, url string, ts string) error {
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if ts == "" {
		req.Header = copyHeader(r.Header, http.Header{}, fwdHeaders)
	} else {
		req.Header = copyHeader(r.Header, http.Header{"Range": {fmt.Sprintf("bytes=%s", ts)}}, []string{"User-Agent"})
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusPartialContent && ts != "" {
		resp.StatusCode = http.StatusOK
		resp.Header.Del("Content-Range")
	}
	to := w.Header()
	copyHeader(resp.Header, to, exposeHeaders)
	to.Set("Cache-Control", "public, max-age=604800")
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	return err
}

func copyHeader(from http.Header, to http.Header, headers []string) http.Header {
	for _, k := range headers {
		if v := from.Get(k); v != "" {
			to.Set(k, v)
		}
	}
	return to
}
