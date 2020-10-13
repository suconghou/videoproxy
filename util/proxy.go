package util

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// MakeClient based proxy
func MakeClient(conf string, timeout time.Duration) http.Client {
	var (
		transport *http.Transport
		err       error
	)
	if strings.HasPrefix(conf, "http") {
		if transport, err = MakeHTTPProxy(conf); err != nil {
			return http.Client{Timeout: timeout}
		}
	} else {
		if transport, err = MakeSocksProxy(conf); err != nil {
			return http.Client{Timeout: timeout}
		}
	}
	return http.Client{Transport: transport, Timeout: timeout}
}

// MakeSocksProxy return socks proxy Transport
func MakeSocksProxy(addr string) (*http.Transport, error) {
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	return &http.Transport{Dial: dialer.Dial, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, nil
}

// MakeHTTPProxy return http proxy Transport
func MakeHTTPProxy(addr string) (*http.Transport, error) {
	urlproxy, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return &http.Transport{Proxy: http.ProxyURL(urlproxy), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, nil
}
