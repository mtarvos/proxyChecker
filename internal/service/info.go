package service

import (
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
	GetInfo(ip string) (entity.IPInfo, error)
}

type InfoRep interface {
	UpdateProxyInfo(proxyItem entity.ProxyItem) error
	GetProxyListForInfo() ([]entity.ProxyItem, error)
}

func NewInfoService(log *slog.Logger, infoProvider InfoProvider, infoRep InfoRep) *InfoService {
	return &InfoService{log: log, infoProvider: infoProvider, infoRep: infoRep}
}

func (i *InfoService) StartInfoRoutine(routineCount int) {
	var wg sync.WaitGroup

	forInfo := make(chan entity.ProxyItem)
	forUpdateInfo := make(chan entity.ProxyItem)

	wg.Add(2)
	go i.updateProxyInfo(forUpdateInfo, &wg)
	go i.fetcherProxyRoutine(forInfo, &wg)

	for q := 0; q < routineCount; q++ {
		wg.Add(1)
		go i.infoRoutine(forInfo, forUpdateInfo, &wg)
	}

	wg.Wait()
}

func (i *InfoService) infoRoutine(forInfo <-chan entity.ProxyItem, forUpdateInfo chan<- entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.infoRoutine"

	for proxyItem := range forInfo {
		i.log.Info("Get info", slog.String("ip", proxyItem.OutIP.String))

		info, err := i.infoProvider.GetInfo(proxyItem.OutIP.String)
		if err != nil {
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
}

func (i *InfoService) updateProxyInfo(forUpdateInfo <-chan entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.updateProxyInfo"

	for proxyItem := range forUpdateInfo {
		if err := i.infoRep.UpdateProxyInfo(proxyItem); err != nil {
			i.log.Error("can not update info for ip", slog.String("ip", proxyItem.IP), slog.String("fn", fn), slog.String("error", err.Error()))
		}
	}

}

func (i *InfoService) fetcherProxyRoutine(forInfo chan<- entity.ProxyItem, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.fetcherProxyRoutine"

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proxyList, err := i.infoRep.GetProxyListForInfo()
			if err != nil {
				i.log.Error("can not get proxy list for get info", slog.String("fn", fn), slog.String("error", err.Error()))
				continue
			}

			for _, item := range proxyList {
				forInfo <- item
			}
		}
	}
}
