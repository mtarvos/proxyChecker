package handler

import (
	"log/slog"
	"net/http"
	"proxyChecker/internal/controller/http/middleware"
)

func (h *Handler) Next() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.next"

		log := h.log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("fn", fn),
		)

		filter, err := h.parseQueryParamsToFilter(r)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter.Page = 0
		filter.Format = "json"

		proxyList, err := h.nextService.GetNextProxy(h.ctx, filter)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.prepareResultWithFormat(w, filter, proxyList)
	}
}
