package handler

import (
	"context"
	"log/slog"
	"proxyChecker/internal/entity"
)

type Handler struct {
	log          *slog.Logger
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

func NewHandler(log *slog.Logger, proxyService ProxyService, statsService StatsService, nextService NextService) *Handler {
	return &Handler{log: log, proxyService: proxyService, statsService: statsService, nextService: nextService}
}

type NextService interface {
	GetNextProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error)
}
