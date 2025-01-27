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

func (s *Storage) UpdateProxyInfo(info entity.IPInfo) error {
	var proxyItem entity.ProxyItem
	proxyItem.OutIP = info.IP
	proxyItem.Country = info.Country
	proxyItem.City = info.City
	proxyItem.ISP = info.ISP
	proxyItem.Timezone = info.Timezone
	return s.UpdateProxyItemByOutIP(proxyItem)
}
