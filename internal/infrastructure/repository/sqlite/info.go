package sqlite

import (
	"database/sql"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

func (s *Storage) GetProxyListForInfo() ([]entity.ProxyItem, error) {
	Alive := 2
	return s.GetProxy(entity.Filters{
		OutIP:   helpers.Cf(sql.NullString{Valid: false}, entity.Ne),
		Country: helpers.Cf(sql.NullString{Valid: false}, entity.Eq),
		Alive:   &Alive,
	})
}

func (s *Storage) UpdateProxyInfo(proxyItem entity.ProxyItem) error {
	return s.UpdateProxyItemByID(proxyItem)
}
