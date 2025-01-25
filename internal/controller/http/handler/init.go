package handler

import (
	"log/slog"
	"proxyChecker/internal/entity"
)

type Handler struct {
	log          *slog.Logger
	proxyService ProxyService
	statsService StatsService
}

type ProxyService interface {
	GetProxyList(filter entity.Filters) ([]entity.ProxyItem, error)
}

type StatsService interface {
	GetStats() (entity.StatsData, error)
}

func NewHandler(log *slog.Logger, proxyService ProxyService, statsService StatsService) *Handler {
	return &Handler{log: log, proxyService: proxyService, statsService: statsService}
}
