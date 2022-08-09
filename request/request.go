package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
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
	}
	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 32*1024))
		},
	}
	errTimeout = errors.New("timeout")
)

type cacheItem struct {
	time    int64
	ctx     context.Context
	cancel  context.CancelFunc
	data    *bytes.Buffer
	headers http.Header
	status  int
	err     error
	loading bool
}

// LockGeter for http cache & lock get
type LockGeter struct {
	time   int64
	caches sync.Map
}

var (
	HttpProvider = NewLockGeter()
)

// NewLockGeter create new lockgeter
func NewLockGeter() *LockGeter {
	return &LockGeter{
		time:   0,
		caches: sync.Map{},
	}
}

func (l *LockGeter) Get(url string, client http.Client, reqHeaders http.Header, ttl int64) ([]byte, http.Header, int, error) {
	var now = time.Now().Unix()
	l.clean(now)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t, loaded := l.caches.LoadOrStore(url, &cacheItem{
		time:    now + ttl,
		ctx:     ctx,
		cancel:  cancel,
		err:     errTimeout,
		loading: true,
	})
	v := t.(*cacheItem)
	if loaded {
		<-v.ctx.Done()
		v.loading = false
		if v.data == nil {
			return nil, v.headers, v.status, v.err
		}
		return v.data.Bytes(), v.headers, v.status, v.err
	}
	data, headers, status, err := Get(url, client, reqHeaders)
	v.data = data
	v.headers = headers
	v.status = status
	v.err = err
	v.loading = false
	cancel()
	if data == nil {
		return nil, headers, status, err
	}
	return data.Bytes(), headers, status, err
}

func (l *LockGeter) clean(now int64) {
	if now-l.time < 5 {
		return
	}
	l.time = now
	l.caches.Range(func(key, value interface{}) bool {
		var v = value.(*cacheItem)
		if v.time < now && !v.loading {
			v.cancel()
			if v.data != nil {
				v.data.Reset()
				bufferPool.Put(v.data)
			}
			l.caches.Delete(key)
		}
		return true
	})
}

// GetByCacher check cache and get from url
func GetByCacher(url string, client http.Client, reqHeaders http.Header) ([]byte, http.Header, int, error) {
	return HttpProvider.Get(url, client, reqHeaders, 86400)
}

// Get http data, the return value should be readonly
func Get(url string, client http.Client, reqHeaders http.Header) (*bytes.Buffer, http.Header, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, 0, err
	}
	req.Header = reqHeaders
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, resp.Header, resp.StatusCode, fmt.Errorf("%s : %s", url, resp.Status)
	}
	var buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		buffer.Reset()
		bufferPool.Put(buffer)
		return nil, resp.Header, resp.StatusCode, err
	}
	return buffer, resp.Header, resp.StatusCode, nil
}

// ProxyData only do get request and pipe without range
func ProxyData(w http.ResponseWriter, r *http.Request, url string, client http.Client) error {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
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
	to.Set("Access-Control-Allow-Origin", "*")
	to.Set("Access-Control-Max-Age", "864000")
	if rhead := r.Header.Get("Access-Control-Request-Headers"); rhead != "" {
		to.Set("Access-Control-Allow-Headers", rhead)
	}
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusPartialContent {
		to.Set("Cache-Control", "public, max-age=864000")
	}
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	return err
}

// Pipe Proxy get request full featured with cache-control & range
func Pipe(w http.ResponseWriter, r *http.Request, url string, client http.Client, rewriteHeader func(http.Header, http.Header)) error {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
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
	to.Set("Access-Control-Allow-Origin", "*")
	to.Set("Access-Control-Max-Age", "864000")
	if rhead := r.Header.Get("Access-Control-Request-Headers"); rhead != "" {
		to.Set("Access-Control-Allow-Headers", rhead)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent || resp.StatusCode == http.StatusNotModified {
		to.Set("Cache-Control", "public, max-age=864000")
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
func ProxyCall(w http.ResponseWriter, url string, client http.Client, rh http.Header, hook func([]byte, int)) error {
	bs, outHeaders, status, err := GetByCacher(url, client, copyHeader(rh, http.Header{}, fwdHeadersBasic))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	var h = w.Header()
	copyHeader(outHeaders, h, exposeHeadersBasic)
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Max-Age", "864000")
	if status == http.StatusOK {
		h.Set("Cache-Control", "public,max-age=864000")
	}
	w.WriteHeader(status)
	_, err = w.Write(bs)
	if hook != nil {
		hook(bs, status)
	}
	return err
}
