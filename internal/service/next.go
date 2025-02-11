package service

import (
	"context"
	"proxyChecker/internal/entity"
)

type NextService struct {
	nextRepo nextRepository
}

type nextRepository interface {
	GetNextProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error)
}

func NewNextService(storage nextRepository) *NextService {
	return &NextService{nextRepo: storage}
}

func (p *NextService) GetNextProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error) {
	return p.nextRepo.GetNextProxy(ctx, filter)
}
