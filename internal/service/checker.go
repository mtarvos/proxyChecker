package service

import (
	"context"
	"errors"
	"log/slog"
	"proxyChecker/internal/entity"
	"sync"
	"time"
)

type CheckerService struct {
	checkerClient CheckerClient
	storage       ProxyStorage
	log           *slog.Logger
}

type aliveStatus struct {
	ProxyAddr
	ip    string
	alive bool
}

type ProxyAddr struct {
	proxyIP   string
	proxyPORT int
}

type CheckerClient interface {
	Check(ctx context.Context, proxyItem entity.ProxyItem) (string, error)
}

type ProxyStorage interface {
	SetAlive(ctx context.Context, proxyItem entity.ProxyItem) error
	GetProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error)
}

func NewCheckerService(log *slog.Logger, storage ProxyStorage, checkerClient CheckerClient) *CheckerService {
	return &CheckerService{log: log, storage: storage, checkerClient: checkerClient}
}

func (c *CheckerService) StartCheckerRoutine(ctx context.Context, routineCount int, checkerWG *sync.WaitGroup) {
	defer checkerWG.Done()

	var wg sync.WaitGroup

	forCheck := make(chan entity.ProxyItem)
	forSetAlive := make(chan entity.ProxyItem)

	wg.Add(2)
	go c.setProxyAliveStatus(ctx, forSetAlive, &wg)
	go c.fetcherProxyRoutine(ctx, forCheck, &wg)

	for i := 0; i < routineCount; i++ {
		wg.Add(1)
		go c.checkerRoutine(ctx, forCheck, forSetAlive, &wg)
	}

	wg.Wait()
	c.log.Info("All Checker goroutine completed")
}

func (c *CheckerService) checkerRoutine(ctx context.Context, forCheckAlive <-chan entity.ProxyItem, forSetAlive chan<- entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.checkerRoutine"
	for proxyItem := range forCheckAlive {
		c.log.Info("Check proxy", slog.String("ip", proxyItem.IP), slog.Int("port", proxyItem.Port))

		outIP, err := c.checkerClient.Check(ctx, proxyItem)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				continue
			}
			c.log.Warn(
				"can not check proxy or proxy is dead!",
				slog.String("fn", fn),
				slog.String("error", err.Error()),
				slog.String("proxyIP", proxyItem.IP),
				slog.Int("proxyPORT", proxyItem.Port),
			)
			proxyItem.Alive.Scan(1)
		} else {
			if proxyItem.OutIP.String != outIP {
				proxyItem.Country.Valid = false
				proxyItem.City.Valid = false
				proxyItem.ISP.Valid = false
				proxyItem.Timezone.Valid = false
			}
			proxyItem.OutIP.Scan(outIP)
			proxyItem.Alive.Scan(2)

			c.log.Info(
				"Proxy is alive!",
				slog.String("ip", proxyItem.IP),
				slog.Int("port", proxyItem.Port),
				slog.String("Out IP", proxyItem.OutIP.String),
			)
		}

		forSetAlive <- proxyItem
	}

	close(forSetAlive)
}

func (c *CheckerService) setProxyAliveStatus(ctx context.Context, forSetAlive <-chan entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.saverProxyAlive"

	for proxyItem := range forSetAlive {
		if err := c.storage.SetAlive(ctx, proxyItem); err != nil {
			if errors.Is(err, context.Canceled) {
				continue
			}
			c.log.Error("can not set Alive for proxy", slog.String("fn", fn), slog.String("error", err.Error()))
		}
	}
}

func (c *CheckerService) fetcherProxyRoutine(ctx context.Context, toCheckerRoutine chan<- entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.fetcherProxyRoutine"

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proxyList, err := c.storage.GetProxy(ctx, entity.Filters{})
			if err != nil {
				if errors.Is(err, context.Canceled) {
					continue
				}
				c.log.Error("can not get proxy list for checking", slog.String("fn", fn), slog.String("error", err.Error()))
				continue
			}

			for _, item := range proxyList {
				if ctx.Err() != nil {
					break
				}
				toCheckerRoutine <- item
			}
		case <-ctx.Done():
			c.log.Info("Context cancelled, stopping fetcherProxyRoutine processing")
			close(toCheckerRoutine)
			return
		}
	}
}
