package handler

import (
	"fmt"
	"net/http"
	"proxyChecker/internal/entity"
	"proxyChecker/internal/lib/helpers"
	"strings"
)

const (
	statusElite        = "ELITE"
	statusNonAnonymous = "TRANSPARENT"
)

func proxyCheckerResponse(ip string, info string, status string) entity.ProxyCheckerResponse {
	return entity.ProxyCheckerResponse{Status: status, Info: info, IP: ip}
}

func (h *Handler) Check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.check"

		ip := strings.Split(r.RemoteAddr, ":")[0]

		err := isProxyHeader(r.Header)
		if err != nil {
			helpers.JSON(w, proxyCheckerResponse(ip, err.Error(), statusNonAnonymous), http.StatusBadRequest)
			return
		}

		helpers.JSON(w, proxyCheckerResponse(ip, "", statusElite), http.StatusOK)
		return
	}
}

func isProxyHeader(headers http.Header) error {
	proxyIndicators := []string{
		"x-forwarded-for",
		"x-forwarded-proto",
		"x-forwarded-host",
		"x-forwarded-port",
		"via",
		"forwarded",
		"x-real-ip",
	}

	for _, indicator := range proxyIndicators {
		ind := headers.Get(indicator)
		if ind != "" {
			return fmt.Errorf("the header have a proxy field: %s", indicator)
		}
	}

	return nil
}
