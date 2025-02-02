package helpers

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func SendGetRequest(url string) (string, error) {
	const fn = "client.sendQuery"

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("%s error send query %s", fn, err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s query return status code: %d", fn, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%s error read body %s", fn, err.Error())
	}

	return string(body), nil
}

func SendGetRequestThroughSocks(ip string, port int, url string) (int, string, error) {
	const fn = "client.SendGetRequestThroughSocks"

	proxyAddr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return 0, "", fmt.Errorf("%s error create dialer %s", fn, err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial:              dialer.Dial,
			DisableKeepAlives: true,
		},
		Timeout: 60 * time.Second,
	}

	return SendRequestThroughClient(client, url)
}

func SendGetRequestThroughHttpProxy(ip string, port int, URL string) (int, string, error) {
	const fn = "client.SendGetRequestThroughHttpProxy"

	proxyAddr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return 0, "", fmt.Errorf("%s bad URL %s: %s", fn, proxyAddr, err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
		Timeout: 5 * time.Second,
	}

	return SendRequestThroughClient(client, URL)
}

func SendGetRequestThroughHttpsProxy(ip string, port int, URL string) (int, string, error) {
	const fn = "client.SendGetRequestThroughHttpsProxy"

	proxyAddr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	proxyURL, err := url.Parse("https://" + proxyAddr)
	if err != nil {
		return 0, "", fmt.Errorf("%s bad url %s: %s", fn, proxyAddr, err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
		Timeout: 5 * time.Second,
	}

	return SendRequestThroughClient(client, URL)
}

func SendRequestThroughClient(client *http.Client, URL string) (int, string, error) {
	const fn = "client.SendRequestThroughClient"

	res, err := client.Get(URL)
	if err != nil {
		return 0, "", fmt.Errorf("%s error send query %s", fn, err.Error())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("%s query return status code: %d", fn, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", fmt.Errorf("%s error read body %s", fn, err.Error())
	}

	return res.StatusCode, string(body), nil
}
