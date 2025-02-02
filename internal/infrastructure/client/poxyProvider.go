package client

import (
	"fmt"
	"log/slog"
	"net"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strconv"
	"strings"
)

type ProxyProvider struct {
	log *slog.Logger
}

func NewProxyProvider(log *slog.Logger) *ProxyProvider {
	return &ProxyProvider{log: log}
}

func (u *ProxyProvider) GetProxies(url string) ([]entity.ProxyItem, error) {
	const fn = "client.GetProxies"

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
	const fn = "client.prepareProxyList"

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

		proxyList = append(proxyList, entity.ProxyItem{IP: ip, Port: iPort})
	}

	return proxyList, nil
}
