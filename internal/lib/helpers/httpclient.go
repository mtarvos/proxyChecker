package helpers

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func SendGetRequest(ctx context.Context, url string) (int, string, error) {
	const fn = "client.SendGetRequest"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	return SendRequestThroughClient(ctx, client, url)
}

func SendGetRequestThroughSocks(ctx context.Context, ip string, port int, url string) (int, string, error) {
	const fn = "client.SendGetRequestThroughSocks"

	proxyAddr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return 0, "", fmt.Errorf("%s error create dialer: %w", fn, err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial:              dialer.Dial,
			DisableKeepAlives: true,
		},
		Timeout: 60 * time.Second,
	}

	return SendRequestThroughClient(ctx, client, url)
}

func SendGetRequestThroughHttpProxy(ctx context.Context, ip string, port int, URL string) (int, string, error) {
	const fn = "client.SendGetRequestThroughHttpProxy"

	proxyAddr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return 0, "", fmt.Errorf("%s bad URL %s: %w", fn, proxyAddr, err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
		Timeout: 5 * time.Second,
	}

	return SendRequestThroughClient(ctx, client, URL)
}

func SendGetRequestThroughHttpsProxy(ctx context.Context, ip string, port int, URL string) (int, string, error) {
	const fn = "client.SendGetRequestThroughHttpsProxy"

	proxyAddr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	proxyURL, err := url.Parse("https://" + proxyAddr)
	if err != nil {
		return 0, "", fmt.Errorf("%s bad url %s: %w", fn, proxyAddr, err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
		Timeout: 5 * time.Second,
	}

	return SendRequestThroughClient(ctx, client, URL)
}

func SendRequestThroughClient(ctx context.Context, client *http.Client, URL string) (int, string, error) {
	const fn = "client.SendRequestThroughClient"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return 0, "", fmt.Errorf("%s error creating request: %w", fn, err)
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("%s error send query: %w", fn, err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("%s query return status code: %d", fn, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", fmt.Errorf("%s error read body: %w", fn, err)
	}

	return res.StatusCode, string(body), nil
}
