package util

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// MakeClient based proxy
func MakeClient(key string, timeout time.Duration) http.Client {
	var (
		conf      = os.Getenv(key)
		transport *http.Transport
		err       error
	)
	if strings.HasPrefix(conf, "http") {
		if transport, err = MakeHTTPProxy(conf); err != nil {
			return http.Client{Timeout: timeout}
		}
	} else if len(conf) > 1 {
		if transport, err = MakeSocksProxy(conf, os.Getenv(key+"_USER"), os.Getenv(key+"_PASSWORD")); err != nil {
			return http.Client{Timeout: timeout}
		}
	}
	return http.Client{Transport: transport, Timeout: timeout}
}

// MakeSocksProxy return socks proxy Transport
func MakeSocksProxy(addr string, user string, password string) (*http.Transport, error) {
	var (
		dialer proxy.Dialer
		err    error
	)
	if user != "" && password != "" {
		dialer, err = proxy.SOCKS5("tcp", addr, &proxy.Auth{User: user, Password: password}, proxy.Direct)
	} else {
		dialer, err = proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	}
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
