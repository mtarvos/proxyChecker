package service

import (
	"context"
	"errors"
	"log/slog"
	"proxyChecker/internal/entity"
	"sync"
	"time"
)

type InfoService struct {
	infoProvider InfoProvider
	infoRep      InfoRep
	log          *slog.Logger
}

type InfoProvider interface {
	GetInfo(ctx context.Context, ip string) (entity.IPInfo, error)
}

type InfoRep interface {
	UpdateProxyInfo(ctx context.Context, proxyItem entity.ProxyItem) error
	GetProxyListForInfo(ctx context.Context) ([]entity.ProxyItem, error)
}

func NewInfoService(log *slog.Logger, infoProvider InfoProvider, infoRep InfoRep) *InfoService {
	return &InfoService{log: log, infoProvider: infoProvider, infoRep: infoRep}
}

func (i *InfoService) StartInfoRoutine(ctx context.Context, routineCount int, infoWG *sync.WaitGroup) {
	defer infoWG.Done()

	var wg sync.WaitGroup

	forInfo := make(chan entity.ProxyItem)
	forUpdateInfo := make(chan entity.ProxyItem)

	wg.Add(1)
	go i.updateProxyInfo(ctx, forUpdateInfo, &wg)
	wg.Add(1)
	go i.fetcherProxyRoutine(ctx, forInfo, &wg)

	for q := 0; q < routineCount; q++ {
		wg.Add(1)
		go i.infoRoutine(ctx, forInfo, forUpdateInfo, &wg)
	}

	wg.Wait()
	i.log.Info("All info goroutine completed")
}

func (i *InfoService) infoRoutine(ctx context.Context, forInfo <-chan entity.ProxyItem, forUpdateInfo chan<- entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.infoRoutine"

	for proxyItem := range forInfo {
		i.log.Info("Get info", slog.String("ip", proxyItem.OutIP.String))

		info, err := i.infoProvider.GetInfo(ctx, proxyItem.OutIP.String)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				continue
			}
			i.log.Error(
				"Can not get info",
				slog.String("fn", fn),
				slog.String("error", err.Error()),
				slog.String("ip", proxyItem.OutIP.String),
			)
			continue
		}

		proxyItem.Country.Scan(info.Country)
		proxyItem.City.Scan(info.City)
		proxyItem.ISP.Scan(info.ISP)
		proxyItem.Timezone.Scan(info.Timezone)

		i.log.Info(
			"Get info success!",
			slog.String("ip", proxyItem.OutIP.String),
			slog.String("Country", proxyItem.Country.String),
			slog.String("City", proxyItem.City.String),
			slog.String("ISP", proxyItem.ISP.String),
			slog.Int("Timezone", int(proxyItem.Timezone.Int32)),
		)

		forUpdateInfo <- proxyItem
	}
	close(forUpdateInfo)

	//i.log.Info("Context cancelled, stopping infoRoutine processing")
}

func (i *InfoService) updateProxyInfo(ctx context.Context, forUpdateInfo <-chan entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.updateProxyInfo"

	for proxyItem := range forUpdateInfo {
		if err := i.infoRep.UpdateProxyInfo(ctx, proxyItem); err != nil {
			if errors.Is(err, context.Canceled) {
				continue
			}
			i.log.Error("can not update info for ip", slog.String("ip", proxyItem.IP), slog.String("fn", fn), slog.String("error", err.Error()))
		}
	}

	//i.log.Info("Context cancelled, stopping updateProxyInfo processing")
}

func (i *InfoService) fetcherProxyRoutine(ctx context.Context, forInfo chan<- entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.fetcherProxyRoutine"

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proxyList, err := i.infoRep.GetProxyListForInfo(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					continue
				}
				i.log.Error("can not get proxy list for get info", slog.String("fn", fn), slog.String("error", err.Error()))
				continue
			}

			for _, item := range proxyList {
				if ctx.Err() != nil {
					break
				}
				forInfo <- item
			}
		case <-ctx.Done():
			//i.log.Info("Context cancelled, stopping fetcherProxyRoutine processing")
			close(forInfo)
			return
		}
	}
}
