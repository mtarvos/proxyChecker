package sqlite

import (
	"context"
	"proxyChecker/internal/entity"
)

func (s *Storage) SetAlive(ctx context.Context, proxyItem entity.ProxyItem) error {
	return s.UpdateProxyItemByID(ctx, proxyItem)
}
