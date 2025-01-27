package service

import (
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
	GetCountByFilter(filter entity.Filters) (int, error)
	GetDistinctField(fieldName string, filter entity.Filters) ([]string, error)
}

func NewStatsService(log *slog.Logger, statsProvider StatsProvider) *StatsService {
	return &StatsService{log: log, provider: statsProvider}
}

func (s *StatsService) GetStats() (entity.StatsData, error) {
	const fn = "stats.GetStats"

	var statsData entity.StatsData

	var err error
	statsData.Total, err = s.provider.GetCountByFilter(entity.Filters{})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get total count: %s", fn, err.Error())
	}

	val := true
	statsData.Alive, err = s.provider.GetCountByFilter(entity.Filters{AliveOnly: &val})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get alive count: %s", fn, err.Error())
	}

	statsData.Dead = statsData.Total - statsData.Alive

	uniqCountry, err := s.provider.GetDistinctField("country", entity.Filters{})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get distinct country count: %s", fn, err.Error())
	}
	statsData.UniqCountry = len(uniqCountry)

	for _, country := range uniqCountry {
		var item entity.CountryStatsItem
		item.Country = country
		item.Count, err = s.provider.GetCountByFilter(
			entity.Filters{Country: helpers.Cf(country, entity.Eq)},
		)

		statsData.CountryStats = append(statsData.CountryStats, item)
	}

	uniqISP, err := s.provider.GetDistinctField("ISP", entity.Filters{})
	if err != nil {
		return entity.StatsData{}, fmt.Errorf("%s can not get distinct ISP count: %s", fn, err.Error())
	}
	statsData.UniqISP = len(uniqISP)

	for _, isp := range uniqISP {
		var item entity.ISPStatsItem
		item.ISP = isp
		item.Count, err = s.provider.GetCountByFilter(
			entity.Filters{Country: helpers.Cf(isp, entity.Eq)},
		)
		statsData.ISPStats = append(statsData.ISPStats, item)
	}
	return statsData, nil
}
