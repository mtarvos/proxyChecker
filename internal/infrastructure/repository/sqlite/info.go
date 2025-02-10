package sqlite

import (
	"context"
	"database/sql"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

func (s *Storage) GetProxyListForInfo(ctx context.Context) ([]entity.ProxyItem, error) {
	Alive := 2
	return s.GetProxy(ctx, entity.Filters{
		OutIP:   helpers.Cf(sql.NullString{Valid: false}, entity.Ne),
		Country: helpers.Cf(sql.NullString{Valid: false}, entity.Eq),
		Alive:   &Alive,
	})
}

func (s *Storage) UpdateProxyInfo(ctx context.Context, proxyItem entity.ProxyItem) error {
	return s.UpdateProxyItemByID(ctx, proxyItem)
}
