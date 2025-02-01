package sqlite

import (
	"proxyChecker/internal/entity"
)

func (s *Storage) SetAlive(proxyItem entity.ProxyItem) error {
	return s.UpdateProxyItemByID(proxyItem)
}
