package util

import (
	"crypto/tls"
	"net"
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
		conf             = os.Getenv(key)
		defaultTransport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		transport *http.Transport
		err       error
	)
	if strings.HasPrefix(conf, "http") {
		if transport, err = MakeHTTPProxy(conf); err != nil {
			return http.Client{Timeout: timeout, Transport: defaultTransport}
		}
	} else if len(conf) > 1 {
		if transport, err = MakeSocksProxy(conf, os.Getenv(key+"_USER"), os.Getenv(key+"_PASSWORD")); err != nil {
			return http.Client{Timeout: timeout, Transport: defaultTransport}
		}
	}
	if transport != nil {
		return http.Client{Timeout: timeout, Transport: transport}
	}
	return http.Client{Timeout: timeout, Transport: defaultTransport}
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
	return &http.Transport{
		Dial:                  dialer.Dial,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, nil
}

// MakeHTTPProxy return http proxy Transport
func MakeHTTPProxy(addr string) (*http.Transport, error) {
	urlproxy, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		Proxy:                 http.ProxyURL(urlproxy),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, nil
}
