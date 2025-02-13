package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

type Checker struct {
	checkerURL string
	proxyType  entity.Status
}

func NewChecker(checkerURL string, proxyType entity.Status) *Checker {
	return &Checker{checkerURL: checkerURL, proxyType: proxyType}
}

func (c *Checker) Check(ctx context.Context, proxyItem entity.ProxyItem) (string, error) {
	const fn = "Checker.Check"

	status, res, err := c.SendRequest(ctx, proxyItem.IP, proxyItem.Port, c.checkerURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if status != http.StatusOK {
		return "", nil
	}

	var chkResp entity.ProxyCheckerResponse
	if err = json.Unmarshal([]byte(res), &chkResp); err != nil {
		return "", fmt.Errorf("%s: bad json: %s %w", fn, res, err)
	}

	return chkResp.IP, nil
}

func (c *Checker) SendRequest(ctx context.Context, proxyIP string, proxyPort int, url string) (int, string, error) {
	const fn = "Checker.SendRequest"

	var result string
	var err error
	var status int

	switch c.proxyType {
	case entity.SOCKS:
		status, result, err = helpers.SendGetRequestThroughSocks(ctx, proxyIP, proxyPort, url)
		if err != nil {
			return 0, "", fmt.Errorf("%s: %w", fn, err)
		}
	case entity.HTTP_PROXY:
		status, result, err = helpers.SendGetRequestThroughHttpProxy(ctx, proxyIP, proxyPort, url)
		if err != nil {
			return 0, "", fmt.Errorf("%s: %w", fn, err)
		}
	case entity.HTTPS_PROXY:
		status, result, err = helpers.SendGetRequestThroughHttpsProxy(ctx, proxyIP, proxyPort, url)
		if err != nil {
			return 0, "", fmt.Errorf("%s: %w", fn, err)
		}
	}

	return status, result, nil
}
