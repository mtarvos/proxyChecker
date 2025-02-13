package service

import (
	"context"
	"fmt"
	"log/slog"
	"proxyChecker/internal/entity"
	"proxyChecker/pkg/logging"
	"sync"
	"time"
)

type UpdaterService struct {
	url      string
	provider ProxyProvider
	saver    ProxySaver
}

type ProxyProvider interface {
	GetProxies(ctx context.Context, url string) ([]entity.ProxyItem, error)
}

type ProxySaver interface {
	SaveProxy(ctx context.Context, proxyList []entity.ProxyItem) error
}

func NewUpdaterService(url string, provider ProxyProvider, saver ProxySaver) *UpdaterService {
	return &UpdaterService{
		url:      url,
		provider: provider,
		saver:    saver,
	}
}

func (u *UpdaterService) StartUpdateProxyRoutine(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	const fn = "proxy.StartUpdateProxyRoutine"
	log := logging.L(ctx).With(
		slog.String("routine", "Updater"),
	)
	ctx = logging.ContextWithLogger(ctx, log)

	log.Debug("call", slog.String("func", fn))

	ticker := time.NewTicker(1 * time.Minute)

	if err := u.worker(ctx, u.url, u.provider, u.saver); err != nil {
		log.Error("error updating proxy list", slog.String("func", fn), slog.String("error", err.Error()))
		return
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := u.worker(ctx, u.url, u.provider, u.saver); err != nil {
					log.Error("error updating proxy list", slog.String("func", fn), slog.String("error", err.Error()))
					break
				}
			case <-ctx.Done():
				break
			}
		}
	}()

	log.Info("Update goroutine stopped")
}

func (u *UpdaterService) worker(ctx context.Context, url string, proxyGetter ProxyProvider, saver ProxySaver) error {
	const fn = "proxy.worker"
	log := logging.L(ctx)

	log.Debug("call", slog.String("func", fn))

	proxyList, err := proxyGetter.GetProxies(ctx, url)
	if err != nil {
		return fmt.Errorf("%s update proxy error: %w", fn, err)
	}

	err = saver.SaveProxy(ctx, proxyList)
	if err != nil {
		return fmt.Errorf("%s save proxy list error: %w", fn, err)
	}

	return nil
}
