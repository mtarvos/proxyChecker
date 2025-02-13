package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"proxyChecker/internal/controller/http/middleware"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strings"
)

func (h *Handler) Stats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.Stats"

		log := h.log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("fn", fn),
		)

		statsData, err := h.statsService.GetStats(r.Context())
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res := prepareStatsResult(statsData)

		helpers.Text(w, res, http.StatusOK)
	}
}

func prepareStatsResult(data entity.StatsData) string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString(fmt.Sprintf("Total: %d\n", data.Total))
	stringBuilder.WriteString(fmt.Sprintf("Alive: %d\n", data.Alive))
	stringBuilder.WriteString(fmt.Sprintf("Dead: %d\n", data.Dead))
	stringBuilder.WriteString(fmt.Sprintf("Uniq Country: %d\n", data.UniqCountry))
	stringBuilder.WriteString(fmt.Sprintf("Uniq ISP: %d\n", data.UniqISP))

	for _, country := range data.CountryStats {
		stringBuilder.WriteString(fmt.Sprintf("%s : %d\n", country.Country, country.Count))
	}

	for _, isp := range data.ISPStats {
		stringBuilder.WriteString(fmt.Sprintf("%s : %d\n", isp.ISP, isp.Count))
	}

	return stringBuilder.String()
}
