package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

type Checker struct {
	log        *slog.Logger
	checkerURL string
	proxyType  entity.Status
}

func NewChecker(log *slog.Logger, checkerURL string, proxyType entity.Status) *Checker {
	return &Checker{log: log, checkerURL: checkerURL, proxyType: proxyType}
}

func (c *Checker) Check(proxyItem entity.ProxyItem) (string, error) {
	const fn = "Checker.Check"

	status, res, err := c.SendRequest(proxyItem.IP, proxyItem.Port, c.checkerURL)
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

func (c *Checker) SendRequest(proxyIP string, proxyPort int, url string) (int, string, error) {
	const fn = "Checker.SendRequest"

	var result string
	var err error
	var status int

	switch c.proxyType {
	case entity.SOCKS:
		status, result, err = helpers.SendGetRequestThroughSocks(proxyIP, proxyPort, url)
		if err != nil {
			return 0, "", fmt.Errorf("%s: %s", fn, err.Error())
		}
	case entity.HTTP_PROXY:
		status, result, err = helpers.SendGetRequestThroughHttpProxy(proxyIP, proxyPort, url)
		if err != nil {
			return 0, "", fmt.Errorf("%s: %s", fn, err.Error())
		}
	case entity.HTTPS_PROXY:
		status, result, err = helpers.SendGetRequestThroughHttpsProxy(proxyIP, proxyPort, url)
		if err != nil {
			return 0, "", fmt.Errorf("%s: %s", fn, err.Error())
		}
	}

	return status, result, nil
}
