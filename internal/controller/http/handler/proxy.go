package handler

import (
	"net/http"
	"proxyChecker/internal/lib/logging"
)

func (h *Handler) Proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.proxy"
		log := logging.L(r.Context())

		filter, err := h.parseQueryParamsToFilter(r)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		proxyList, err := h.proxyService.GetProxyList(r.Context(), filter)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.prepareResultWithFormat(r.Context(), w, filter, proxyList)
	}
}
