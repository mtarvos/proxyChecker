package service

import (
	"context"
	"fmt"
	"log/slog"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
)

type StatsService struct {
	log      *slog.Logger
	provider StatsProvider
}

type StatsProvider interface {
	GetCountByFilter(ctx context.Context, filter entity.Filters) (int, error)
	GetDistinctField(ctx context.Context, fieldName string, filter entity.Filters) ([]string, error)
}

func NewStatsService(log *slog.Logger, statsProvider StatsProvider) *StatsService {
	return &StatsService{log: log, provider: statsProvider}
}

func (s *StatsService) GetStats(ctx context.Context) (entity.StatsData, error) {
	const fn = "stats.GetStats"

	var statsData entity.StatsData

	var err error
	statsData.Total, err = s.provider.GetCountByFilter(ctx, entity.Filters{})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get total count: %w", fn, err)
	}

	val := 2
	statsData.Alive, err = s.provider.GetCountByFilter(ctx, entity.Filters{Alive: &val})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get alive count: %w", fn, err)
	}

	statsData.Dead = statsData.Total - statsData.Alive

	uniqCountry, err := s.provider.GetDistinctField(ctx, "country", entity.Filters{})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get distinct country count: %w", fn, err)
	}
	statsData.UniqCountry = len(uniqCountry)

	for _, country := range uniqCountry {
		var item entity.CountryStatsItem
		item.Country = country
		item.Count, err = s.provider.GetCountByFilter(
			ctx,
			entity.Filters{Country: helpers.Cf(country, entity.Eq)},
		)

		statsData.CountryStats = append(statsData.CountryStats, item)
	}

	uniqISP, err := s.provider.GetDistinctField(ctx, "ISP", entity.Filters{})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get distinct ISP count: %w", fn, err)
	}
	statsData.UniqISP = len(uniqISP)

	for _, isp := range uniqISP {
		var item entity.ISPStatsItem
		item.ISP = isp
		item.Count, err = s.provider.GetCountByFilter(
			ctx,
			entity.Filters{Country: helpers.Cf(isp, entity.Eq)},
		)
		statsData.ISPStats = append(statsData.ISPStats, item)
	}
	return statsData, nil
}
