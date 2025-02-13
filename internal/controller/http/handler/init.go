package handler

import (
	"context"
	"proxyChecker/internal/entity"
)

type Handler struct {
	proxyService ProxyService
	statsService StatsService
	nextService  NextService
}

type ProxyService interface {
	GetProxyList(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error)
	GetTotalCountByFilter(ctx context.Context, filter entity.Filters) (int, error)
}

type StatsService interface {
	GetStats(ctx context.Context) (entity.StatsData, error)
}

func NewHandler(proxyService ProxyService, statsService StatsService, nextService NextService) *Handler {
	return &Handler{proxyService: proxyService, statsService: statsService, nextService: nextService}
}

type NextService interface {
	GetNextProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error)
}
