package service

import (
	"proxyChecker/internal/entity"
)

type ProxyService struct {
	proxyRepo proxyRepository
}

type proxyRepository interface {
	Get(filter entity.Filters) ([]entity.ProxyItem, error)
}

func NewProxy(storage proxyRepository) *ProxyService {
	return &ProxyService{proxyRepo: storage}
}

func (p *ProxyService) GetProxyList(filter entity.Filters) ([]entity.ProxyItem, error) {
	return p.proxyRepo.Get(filter)
}
