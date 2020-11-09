package request

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/suconghou/videoproxy/util"
)

var (
	disableGzip = os.Getenv("DISABLE_GZIP")
	fwdHeaders  = []string{
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

	fwdHeadersBasic = []string{
		"User-Agent",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
	}
	exposeHeadersBasic = []string{
		"Content-Length",
		"Content-Type",
		"Content-Encoding",
		"Date",
	}
)

type cacheItem struct {
	data   []byte
	status int
	age    time.Time
}

type bytecache struct {
	sync.RWMutex
	data map[string]cacheItem
	age  time.Duration
}

var (
	longCacher = &bytecache{
		data: make(map[string]cacheItem),
		age:  time.Hour * 48,
	}
)

func (by *bytecache) geturl(url string, client http.Client) ([]byte, int, error) {
	var bs, status = by.get(url)
	if bs != nil {
		return bs, status, nil
	}
	res, status, err := GetURLBody(url, client)
	if err != nil {
		return nil, status, err
	}
	if len(res) > 0 && status == http.StatusOK {
		by.set(url, res, status)
	}
	return res, status, nil
}

func (by *bytecache) get(key string) ([]byte, int) {
	by.RLock()
	item := by.data[key]
	by.RUnlock()
	if item.age.After(time.Now()) {
		return item.data, item.status
	}
	by.expire()
	return nil, 0
}

func (by *bytecache) set(key string, data []byte, status int) {
	by.Lock()
	by.data[key] = cacheItem{data, status, time.Now().Add(by.age)}
	by.Unlock()
}

func (by *bytecache) expire() {
	t := time.Now()
	by.Lock()
	for key, item := range by.data {
		if item.age.Before(t) {
			delete(by.data, key)
		}
	}
	by.Unlock()
}

// GetURLData check cache and get from url
func GetURLData(url string, client http.Client) ([]byte, int, error) {
	return longCacher.geturl(url, client)
}

// GetURLBody run quick get no cache
func GetURLBody(url string, client http.Client) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	return bs, resp.StatusCode, err
}

// ProxyData only do get request and pipe without range
func ProxyData(w http.ResponseWriter, r *http.Request, url string, client http.Client) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	req.Header = copyHeader(r.Header, http.Header{}, fwdHeadersBasic)
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer res.Body.Close()
	to := w.Header()
	copyHeader(res.Header, to, exposeHeadersBasic)
	to.Set("Cache-Control", "public, max-age=864000")
	to.Set("Access-Control-Allow-Origin", "*")
	if rhead := r.Header.Get("Access-Control-Request-Headers"); rhead != "" {
		to.Set("Access-Control-Allow-Headers", rhead)
	}
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	return err
}

// Pipe Proxy get request full featured with cache-control & range
func Pipe(w http.ResponseWriter, r *http.Request, url string, client http.Client, rewriteHeader func(http.Header, http.Header)) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	req.Header = copyHeader(r.Header, http.Header{}, fwdHeaders)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()
	to := w.Header()
	copyHeader(resp.Header, to, exposeHeaders)
	to.Set("Cache-Control", "public, max-age=864000")
	to.Set("Access-Control-Allow-Origin", "*")
	if rhead := r.Header.Get("Access-Control-Request-Headers"); rhead != "" {
		to.Set("Access-Control-Allow-Headers", rhead)
	}
	if rewriteHeader != nil {
		rewriteHeader(resp.Header, to)
	}
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

// ProxyCall call api with long cache
func ProxyCall(w http.ResponseWriter, url string, client http.Client) error {
	bs, status, err := GetURLData(url, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	var (
		h      = w.Header()
		gziped bool
	)
	bs, gziped = gzipResponse(bs)
	if gziped {
		h.Set("Content-Encoding", "gzip")
	}
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Cache-Control", "public,max-age=864000")
	w.WriteHeader(status)
	_, err = w.Write(bs)
	return err
}

func gzipResponse(bs []byte) ([]byte, bool) {
	if disableGzip != "" || len(bs) < 512 {
		return bs, false
	}
	gz, err := util.GzipEncode(bs)
	if err == nil {
		return gz, true
	}
	util.Log.Println("gz error", err)
	return bs, false
}
