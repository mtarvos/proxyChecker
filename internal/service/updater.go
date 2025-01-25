package service

import (
	"fmt"
	"log/slog"
	"proxyChecker/internal/entity"
	"time"
)

type UpdaterService struct {
	log      *slog.Logger
	url      string
	provider ProxyProvider
	saver    ProxySaver
}

type ProxyProvider interface {
	GetProxies(url string) ([]entity.ProxyItem, error)
}

type ProxySaver interface {
	SaveProxy(proxyList []entity.ProxyItem) error
}

func NewUpdaterService(log *slog.Logger, url string, provider ProxyProvider, saver ProxySaver) *UpdaterService {
	return &UpdaterService{
		log:      log,
		url:      url,
		provider: provider,
		saver:    saver,
	}
}

func (u *UpdaterService) StartUpdateProxyThread() {
	const fn = "proxy.StartUpdateProxyThread"

	u.log.Debug("call", slog.String("func", fn))

	ticker := time.NewTicker(1 * time.Minute)

	quit := make(chan struct{})

	if err := u.worker(u.url, u.provider, u.saver); err != nil {
		u.log.Error("error updating proxy list", slog.String("func", fn), slog.String("error", err.Error()))
		return
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := u.worker(u.url, u.provider, u.saver); err != nil {
					u.log.Error("error updating proxy list; exit goroutine", slog.String("func", fn), slog.String("error", err.Error()))
					return
				}
			case <-quit:
				u.log.Debug("%s: update goroutine stopped", fn)
				return
			}
		}
	}()

}

func (u *UpdaterService) worker(url string, proxyGetter ProxyProvider, saver ProxySaver) error {
	const fn = "proxy.worker"

	u.log.Debug("call", slog.String("func", fn))

	proxyList, err := proxyGetter.GetProxies(url)
	if err != nil {
		return fmt.Errorf("%s update proxy error: %s", fn, err.Error())
	}

	err = saver.SaveProxy(proxyList)
	if err != nil {
		return fmt.Errorf("%s save proxy list error: %s", fn, err.Error())
	}

	return nil
}
