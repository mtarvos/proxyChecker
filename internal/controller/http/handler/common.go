package handler

import (
	"context"
	"fmt"
	htmlTemplate "html/template"
	"net/http"
	"proxyChecker/internal/controller/http/templates"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strconv"
	textTemplate "text/template"

	"proxyChecker/pkg/logging"
)

type PageData struct {
	ProxyList   []entity.ProxyItem
	Pages       []int
	CurrentPage int
	TotalPages  int
	Limit       int
}

func (h *Handler) parseQueryParamsToFilter(r *http.Request) (entity.Filters, error) {
	country := r.URL.Query().Get("country")
	city := r.URL.Query().Get("city")
	isp := r.URL.Query().Get("ISP")
	alive := r.URL.Query().Get("alive")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	format := r.URL.Query().Get("format")
	label := r.URL.Query().Get("label")

	var pAlive *int
	if alive != "" {
		valAlive, err := strconv.Atoi(alive)
		if err != nil {
			return entity.Filters{}, fmt.Errorf("invalid value for field alive, need to be int")
		}
		pAlive = &valAlive
	}

	var valPage int
	var err error
	if page != "" {
		valPage, err = strconv.Atoi(page)
		if err != nil {
			return entity.Filters{}, fmt.Errorf("invalid value for field page, need to be int")
		}
		if valPage < 0 {
			valPage = 1
		}
	} else {
		valPage = 1
	}

	var valLimit int
	if limit != "" {
		valLimit, err = strconv.Atoi(limit)
		if err != nil {
			return entity.Filters{}, fmt.Errorf("invalid value for field limit, need to be int")
		}
		if valLimit < 0 {
			valLimit = 10
		}
	} else {
		valLimit = 10
	}

	var filter entity.Filters
	filter.Alive = pAlive
	filter.Page = valPage
	filter.Limit = valLimit

	if country != "" {
		filter.Country = helpers.Cf(country, entity.Eq)
	}
	if city != "" {
		filter.City = helpers.Cf(city, entity.Eq)
	}
	if isp != "" {
		filter.ISP = helpers.Cf(isp, entity.Eq)
	}

	filter.Format = format
	filter.Label = label

	return filter, nil
}

func (h *Handler) prepareResultWithFormat(ctx context.Context, w http.ResponseWriter, filter entity.Filters, proxyList []entity.ProxyItem) {
	log := logging.L(ctx)
	if filter.Format == "json" {
		helpers.JSON(w, proxyList, http.StatusOK)
		return
	}

	var err error
	if filter.Format == "text" {
		if err = h.ProxyListTEXT(w, templates.TEXTProxyList, proxyList, http.StatusOK); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Error(err.Error())
		}
		return
	}

	totalItems, err := h.proxyService.GetTotalCountByFilter(ctx, filter)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err = h.ProxyListHTML(w, templates.HTMLProxyList, proxyList, filter.Page, filter.Limit, totalItems, http.StatusOK); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
}

func (h *Handler) ProxyListHTML(w http.ResponseWriter, tmplHtml string, proxyList []entity.ProxyItem, currentPage int, limit int, totalItems int, status int) error {
	if currentPage == 0 {
		currentPage = 1
	}

	if limit == 0 {
		limit = 10
	}

	pageData := h.newPageData(proxyList, currentPage, limit, totalItems)

	tmpl, err := htmlTemplate.New("proxyList").
		Funcs(htmlTemplate.FuncMap{
			"add":      func(a, b int) int { return a + b },
			"subtract": func(a, b int) int { return a - b },
		}).
		Parse(tmplHtml)
	if err != nil {
		return fmt.Errorf("html template parsing error: %w", err)
	}

	if err = tmpl.Execute(w, pageData); err != nil {
		return fmt.Errorf("html template execute error: %w", err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	return nil
}

func (h *Handler) newPageData(proxyList []entity.ProxyItem, currentPage int, limit int, totalItems int) *PageData {
	const visiblePages = 10
	totalPages := (totalItems + limit - 1) / limit
	startPage := currentPage - visiblePages/2
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + visiblePages - 1
	if endPage > totalPages {
		endPage = totalPages
		startPage = endPage - visiblePages
		if startPage < 1 {
			startPage = 1
		}
	}

	pages := make([]int, 0, endPage-startPage+1)
	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	return &PageData{ProxyList: proxyList, Pages: pages, CurrentPage: currentPage, TotalPages: totalPages, Limit: limit}
}

func (h *Handler) ProxyListTEXT(w http.ResponseWriter, tmpText string, proxyList []entity.ProxyItem, status int) error {
	tmpl, err := textTemplate.New("proxyList").Parse(tmpText)
	if err != nil {
		return fmt.Errorf("text template parsing error: %w", err)
	}

	if err = tmpl.Execute(w, proxyList); err != nil {
		return fmt.Errorf("text template execute error: %w", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)

	return nil
}
