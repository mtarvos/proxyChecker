package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"proxyChecker/internal/controller/http/middleware"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strconv"
	"strings"
)

func (p *Handler) Proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.proxy"

		log := p.log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("fn", fn),
		)

		log.Debug("call")

		filter, err := parseQueryParamsToFilter(r)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		proxyList, err := p.proxyService.GetProxyList(filter)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res := prepareProxyResult(proxyList)

		helpers.Text(w, res, http.StatusOK)
	}
}

func parseQueryParamsToFilter(r *http.Request) (entity.Filters, error) {
	country := r.URL.Query().Get("country")
	isp := r.URL.Query().Get("ISP")

	const aliveQueryParam = "onlyAlive"

	var bAlive *bool
	alive := r.URL.Query().Get("onlyAlive")
	if alive != "" {
		val, err := strconv.ParseBool(alive)
		if err != nil {
			return entity.Filters{}, fmt.Errorf("invalid value for field %s, need to be boolean", aliveQueryParam)
		}

		bAlive = &val
	}

	return entity.Filters{
		AliveOnly: bAlive,
		Country:   helpers.Cf(country, entity.Eq),
		ISP:       helpers.Cf(isp, entity.Eq),
	}, nil
}

func prepareProxyResult(proxyList []entity.ProxyItem) string {
	if len(proxyList) == 0 {
		return ""
	}

	var builder strings.Builder
	for _, proxy := range proxyList {
		builder.WriteString(fmt.Sprintf("%s:%d; %s; %s\n", proxy.IP, proxy.Port, proxy.Country, proxy.ISP))
	}

	return builder.String()
}
