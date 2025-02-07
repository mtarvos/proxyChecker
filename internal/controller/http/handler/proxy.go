package handler

import (
	"fmt"
	htmlTemplate "html/template"
	"log/slog"
	"net/http"
	"proxyChecker/internal/controller/http/middleware"
	"proxyChecker/internal/controller/http/templates"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strconv"
	"strings"
	textTemplate "text/template"
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

		format := r.URL.Query().Get("format")

		if format == "json" {
			helpers.JSON(w, proxyList, http.StatusOK)
			return
		}

		if format == "text" {
			if err = ProxyListTEXT(w, templates.TEXTProxyList, proxyList, http.StatusOK); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Error(err.Error())
			}
			return
		}

		if err = ProxyListHTML(w, templates.HTMLProxyList, proxyList, http.StatusOK); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
	}
}

func parseQueryParamsToFilter(r *http.Request) (entity.Filters, error) {
	country := r.URL.Query().Get("country")
	city := r.URL.Query().Get("city")
	isp := r.URL.Query().Get("ISP")
	alive := r.URL.Query().Get("alive")

	var pAlive *int

	if alive != "" {
		val, err := strconv.Atoi(alive)
		if err != nil {
			return entity.Filters{}, fmt.Errorf("invalid value for field alive, need to be int")
		}

		pAlive = &val
	}

	var filter entity.Filters
	filter.Alive = pAlive
	if country != "" {
		filter.Country = helpers.Cf(country, entity.Eq)
	}
	if city != "" {
		filter.City = helpers.Cf(city, entity.Eq)
	}
	if isp != "" {
		filter.ISP = helpers.Cf(isp, entity.Eq)
	}

	return filter, nil
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

func ProxyListHTML(w http.ResponseWriter, tmplHtml string, proxyList []entity.ProxyItem, status int) error {

	tmpl, err := htmlTemplate.New("proxyList").Parse(tmplHtml)
	if err != nil {
		return fmt.Errorf("html template parsing error: %w", err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	if err = tmpl.Execute(w, proxyList); err != nil {
		return fmt.Errorf("html template execute error: %w", err)
	}
	return nil
}

func ProxyListTEXT(w http.ResponseWriter, tmpText string, proxyList []entity.ProxyItem, status int) error {

	tmpl, err := textTemplate.New("proxyList").Parse(tmpText)
	if err != nil {
		return fmt.Errorf("text template parsing error: %w", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)

	if err = tmpl.Execute(w, proxyList); err != nil {
		return fmt.Errorf("text template execute error: %w", err)
	}
	return nil
}
