package sqlite

import (
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

func (s *Storage) GetProxyListForInfo() ([]entity.ProxyItem, error) {
	return s.GetProxy(entity.Filters{
		OutIP:   helpers.Cf("", entity.Ne),
		Country: helpers.Cf("", entity.Eq),
	})
}

func (s *Storage) UpdateProxyInfo(proxyItem entity.ProxyItem) error {
	return s.UpdateProxyItemByID(proxyItem)
}
