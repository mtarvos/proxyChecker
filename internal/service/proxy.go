package service

import (
	"context"
	"proxyChecker/internal/entity"
)

type ProxyService struct {
	proxyRepo proxyRepository
}

type proxyRepository interface {
	GetProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error)
	GetCountByFilter(ctx context.Context, filter entity.Filters) (int, error)
}

func NewProxy(storage proxyRepository) *ProxyService {
	return &ProxyService{proxyRepo: storage}
}

func (p *ProxyService) GetProxyList(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error) {
	return p.proxyRepo.GetProxy(ctx, filter)
}

func (p *ProxyService) GetTotalCountByFilter(ctx context.Context, filter entity.Filters) (int, error) {
	return p.proxyRepo.GetCountByFilter(ctx, filter)
}
