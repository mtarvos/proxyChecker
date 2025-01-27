package service

import (
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
	Check(ip string, port int) (string, bool, error)
}

type ProxyStorage interface {
	SetAlive(proxy string, port int, ip string, alive bool) error
	GetProxy(filter entity.Filters) ([]entity.ProxyItem, error)
}

func NewCheckerService(log *slog.Logger, storage ProxyStorage, checkerClient CheckerClient) *CheckerService {
	return &CheckerService{log: log, storage: storage, checkerClient: checkerClient}
}

func (c *CheckerService) StartCheckerRoutine(routineCount int) {
	var wg sync.WaitGroup

	forCheck := make(chan ProxyAddr)
	forSetAlive := make(chan aliveStatus)

	wg.Add(2)
	go c.setProxyAliveStatus(forSetAlive, &wg)
	go c.fetcherProxyRoutine(forCheck, &wg)

	for i := 0; i < routineCount; i++ {
		wg.Add(1)
		go c.checkerRoutine(forCheck, forSetAlive, &wg)
	}

	wg.Wait()
}

func (c *CheckerService) checkerRoutine(forCheckAlive <-chan ProxyAddr, forSetAlive chan<- aliveStatus, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.checkerRoutine"

	for proxyAddr := range forCheckAlive {
		outIP, alive, err := c.checkerClient.Check(proxyAddr.proxyIP, proxyAddr.proxyPORT)
		if err != nil {
			c.log.Error(
				"can not check proxy",
				slog.String("fn", fn),
				slog.String("error", err.Error()),
				slog.String("proxyIP", proxyAddr.proxyIP),
				slog.Int("proxyPORT", proxyAddr.proxyPORT),
			)
			continue
		}

		forSetAlive <- aliveStatus{
			ProxyAddr: ProxyAddr{
				proxyIP:   proxyAddr.proxyIP,
				proxyPORT: proxyAddr.proxyPORT,
			},
			ip:    outIP,
			alive: alive,
		}
	}

}

func (c *CheckerService) setProxyAliveStatus(forSetAlive <-chan aliveStatus, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.saverProxyAlive"

	for aliveStatus := range forSetAlive {
		if err := c.storage.SetAlive(aliveStatus.proxyIP, aliveStatus.proxyPORT, aliveStatus.ip, aliveStatus.alive); err != nil {
			c.log.Error("can not set Alive for proxy", slog.String("fn", fn), slog.String("error", err.Error()))
		}
	}

}

func (c *CheckerService) fetcherProxyRoutine(toCheckerRoutine chan<- ProxyAddr, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.fetcherProxyRoutine"

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proxyList, err := c.storage.GetProxy(entity.Filters{})
			if err != nil {
				c.log.Error("can not get proxy list for checking", slog.String("fn", fn), slog.String("error", err.Error()))
				continue
			}

			for _, item := range proxyList {
				toCheckerRoutine <- ProxyAddr{proxyIP: item.IP, proxyPORT: item.Port}
			}
		}
	}
}
