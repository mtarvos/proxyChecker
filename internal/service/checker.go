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
	ip    string
	alive bool
}

type ProxyAddr struct {
	ip   string
	port int
}

type CheckerClient interface {
	Check(ip string, port int) (string, bool, error)
}

type ProxyStorage interface {
	SetAlive(ip string, alive bool) error
	Get(filter entity.Filters) ([]entity.ProxyItem, error)
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
		outIP, alive, err := c.checkerClient.Check(proxyAddr.ip, proxyAddr.port)
		if err != nil {
			c.log.Error(
				"can not check proxy",
				slog.String("fn", fn),
				slog.String("error", err.Error()),
				slog.String("ip", proxyAddr.ip),
				slog.Int("port", proxyAddr.port),
			)
			continue
		}

		forSetAlive <- aliveStatus{ip: outIP, alive: alive}
	}

}

func (c *CheckerService) setProxyAliveStatus(forSetAlive <-chan aliveStatus, wg *sync.WaitGroup) {
	defer wg.Done()

	const fn = "CheckerService.saverProxyAlive"

	for aliveStatus := range forSetAlive {
		if err := c.storage.SetAlive(aliveStatus.ip, aliveStatus.alive); err != nil {
			c.log.Error("can not get proxy list for checking", slog.String("fn", fn), slog.String("error", err.Error()))
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
			proxyList, err := c.storage.Get(entity.Filters{})
			if err != nil {
				c.log.Error("can not get proxy list for checking", slog.String("fn", fn), slog.String("error", err.Error()))
				continue
			}

			for _, item := range proxyList {
				toCheckerRoutine <- ProxyAddr{ip: item.Ip, port: item.Port}
			}
		}
	}
}
