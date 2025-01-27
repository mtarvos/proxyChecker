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
	UpdateProxyInfo(entity.IPInfo) error
	GetProxyListForInfo() ([]entity.ProxyItem, error)
}

func NewInfoService(log *slog.Logger, infoProvider InfoProvider, infoRep InfoRep) *InfoService {
	return &InfoService{log: log, infoProvider: infoProvider, infoRep: infoRep}
}

func (i *InfoService) StartInfoRoutine(routineCount int) {
	var wg sync.WaitGroup

	forInfo := make(chan string)
	forUpdateInfo := make(chan entity.IPInfo)

	wg.Add(2)
	go i.updateProxyInfo(forUpdateInfo, &wg)
	go i.fetcherProxyRoutine(forInfo, &wg)

	for q := 0; q < routineCount; q++ {
		wg.Add(1)
		go i.infoRoutine(forInfo, forUpdateInfo, &wg)
	}

	wg.Wait()
}

func (i *InfoService) infoRoutine(forInfo <-chan string, forUpdateInfo chan<- entity.IPInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.infoRoutine"

	for ip := range forInfo {
		info, err := i.infoProvider.GetInfo(ip)
		if err != nil {
			i.log.Error(
				"can not get info",
				slog.String("fn", fn),
				slog.String("error", err.Error()),
				slog.String("ip", ip),
			)
			continue
		}

		forUpdateInfo <- info
	}
}

func (i *InfoService) updateProxyInfo(forUpdateInfo <-chan entity.IPInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "InfoService.updateProxyInfo"

	for ipInfo := range forUpdateInfo {
		if err := i.infoRep.UpdateProxyInfo(ipInfo); err != nil {
			i.log.Error("can not update info for ip", slog.String("ip", ipInfo.IP), slog.String("fn", fn), slog.String("error", err.Error()))
		}
	}

}

func (i *InfoService) fetcherProxyRoutine(forInfo chan<- string, wg *sync.WaitGroup) {
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
				forInfo <- item.OutIP
			}
		}
	}
}
