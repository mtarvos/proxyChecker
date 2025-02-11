package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"sync"
)

func (s *Storage) GetNextProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error) {
	const fn = "sqlite.GetNextProxy"

	limit := filter.Limit

	Alive := 2
	proxyList, err := s.GetProxy(ctx, entity.Filters{
		OutIP:   helpers.Cf(sql.NullString{Valid: false}, entity.Ne),
		Country: helpers.Cf(sql.NullString{Valid: false}, entity.Ne),
		Alive:   &Alive,
	})
	if err != nil {
		return nil, fmt.Errorf("%s opening sqlite db: %w", fn, err)
	}

	nextProxyList, err := s.getUniqProxyList(proxyList, filter.Label, limit)
	if err != nil {
		return []entity.ProxyItem{}, err
	}

	return nextProxyList, nil
}

func (s *Storage) getUniqProxyList(proxyList []entity.ProxyItem, label string, limit int) ([]entity.ProxyItem, error) {
	var nextProxyList []entity.ProxyItem

	if len(proxyList) == 0 {
		return []entity.ProxyItem{}, nil
	}

	if label == "" {
		label = "defaultLabel"
	}

	for i := 0; i < limit; i++ {
		proxy, err := s.getNext(proxyList, label)
		if err != nil {
			return []entity.ProxyItem{}, err
		}
		nextProxyList = append(nextProxyList, proxy)
	}

	return nextProxyList, nil
}

func (s *Storage) getNext(proxyList []entity.ProxyItem, label string) (entity.ProxyItem, error) {
	const fn = "sqlite.getNext"

	var rwm sync.RWMutex

	if len(proxyList) == 0 {
		return entity.ProxyItem{}, fmt.Errorf("%s: empty proxyList", fn)
	}

	rwm.Lock()
	if s.used[label] == nil {
		s.used[label] = make(map[string]bool)
	}
	rwm.Unlock()

	var proxyItem *entity.ProxyItem

	rwm.RLock()
	for _, proxy := range proxyList {
		if s.used[label][proxy.OutIP.String] {
			continue
		}
		proxyItem = &proxy
		break
	}
	rwm.RUnlock()

	if proxyItem == nil {
		rwm.Lock()
		delete(s.used, label)
		rwm.Unlock()
		return s.getNext(proxyList, label)
	}

	rwm.Lock()
	s.used[label][proxyItem.OutIP.String] = true
	rwm.Unlock()

	return *proxyItem, nil
}
