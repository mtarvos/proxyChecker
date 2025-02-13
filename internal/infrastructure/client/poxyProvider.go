package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"proxyChecker/internal/lib/logging"
	"strconv"
	"strings"
)

type ProxyProvider struct {
}

func NewProxyProvider() *ProxyProvider {
	return &ProxyProvider{}
}

func (u *ProxyProvider) GetProxies(ctx context.Context, url string) ([]entity.ProxyItem, error) {
	const fn = "client.GetProxies"
	log := logging.L(ctx)

	log.Debug("call", slog.String("url", url))

	_, res, err := helpers.SendGetRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%s SendGetRequest: %w", fn, err)
	}
	log.Debug("response", slog.String("result", res))

	proxyList, err := prepareProxyList(res)
	if err != nil {
		return nil, fmt.Errorf("%s prepareProxyList: %w", fn, err)
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
			return nil, fmt.Errorf("%s error parse ip and port from query result: %w", fn, err)
		}

		iPort, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("%s error convert port to string: %w", fn, err)
		}

		proxyList = append(proxyList, entity.ProxyItem{IP: ip, Port: iPort})
	}

	return proxyList, nil
}
