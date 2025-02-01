package external

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

type CheckerApiClient struct {
	log        *slog.Logger
	checkerURL string
	proxyType  entity.Status
}

func NewCheckerApiClient(log *slog.Logger, checkerURL string, proxyType entity.Status) *CheckerApiClient {
	return &CheckerApiClient{log: log, checkerURL: checkerURL, proxyType: proxyType}
}

func (c *CheckerApiClient) Check(proxyItem entity.ProxyItem) (string, error) {
	const fn = "CheckerApiClient.Check"

	c.log.Debug("call", slog.String("func", fn), slog.String("proxy-ip", proxyItem.IP), slog.Int("port", proxyItem.Port))

	status, res, err := c.SendRequest(proxyItem.IP, proxyItem.Port, c.checkerURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	c.log.Debug("response", slog.String("func", fn), slog.Int("statusCode", status), slog.String("result", res))

	if status != http.StatusOK {
		return "", nil
	}

	var chkResp entity.ProxyCheckerResponse
	if err = json.Unmarshal([]byte(res), &chkResp); err != nil {
		return "", fmt.Errorf("%s: bad json: %s %w", fn, res, err)
	}

	return chkResp.IP, nil
}

func (c *CheckerApiClient) SendRequest(proxyIP string, proxyPort int, url string) (int, string, error) {
	const fn = "CheckerApiClient.SendRequest"

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
