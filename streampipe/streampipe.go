package streampipe

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func Pipe(w http.ResponseWriter, r *http.Request, url string, rewriteHeader func(http.ResponseWriter, int, http.Header)) error {
	client, req, err := initClient(r, url)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return err
	}
	if res, err := client.Do(req); err == nil {
		for key, value := range res.Header {
			w.Header().Set(key, value[0])
		}
		w.WriteHeader(res.StatusCode)
		if rewriteHeader != nil {
			rewriteHeader(w, res.StatusCode, res.Header)
		}
		defer res.Body.Close()
		if _, err := io.Copy(w, res.Body); err == nil {
			return nil
		} else {
			return err
		}

	} else {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return err
	}

}

func Get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	str, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return str, err
	}
	return str, nil
}

func initClient(r *http.Request, url string) (*http.Client, *http.Request, error) {
	cli := &http.Client{Timeout: 3600 * time.Second}
	if req, err := http.NewRequest("GET", url, nil); err == nil {
		return cli, reqWithHeader(req, r), nil
	} else {
		return cli, nil, err
	}
}

func reqWithHeader(req *http.Request, r *http.Request) *http.Request {
	req.Header = r.Header
	return req
}
