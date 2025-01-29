package external

import (
	"fmt"
	"log/slog"
	"net"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strconv"
	"strings"
)

type ProxyApiClient struct {
	log *slog.Logger
}

func NewProxyApiClient(log *slog.Logger) *ProxyApiClient {
	return &ProxyApiClient{log: log}
}

func (u *ProxyApiClient) GetProxies(url string) ([]entity.ProxyItem, error) {
	const fn = "external.GetProxies"

	u.log.Debug("call", slog.String("func", fn), slog.String("url", url))

	res, err := helpers.SendGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", fn, err.Error())
	}

	u.log.Debug("response", slog.String("func", fn), slog.String("result", res))

	proxyList, err := prepareProxyList(res)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return proxyList, nil
}

func prepareProxyList(text string) ([]entity.ProxyItem, error) {
	const fn = "external.prepareProxyList"

	if len(text) == 0 {
		return []entity.ProxyItem{}, nil
	}

	list := strings.Fields(text)

	var proxyList []entity.ProxyItem

	for _, item := range list {
		ip, port, err := net.SplitHostPort(item)
		if err != nil {
			return nil, fmt.Errorf("error parse ip and port from query result", fn, err.Error())
		}

		iPort, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("error convert port to string", fn, err.Error())
		}

		proxyList = append(proxyList, entity.ProxyItem{Ip: ip, Port: iPort})
	}

	return proxyList, nil
}
