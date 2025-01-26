package handler

import (
	"fmt"
	"net/http"
	"proxyChecker/internal/lib/helpers"
	"strings"
)

const (
	statusProxy    = "PROXY"
	statusNonProxy = "NON_PROXY"
)

type ResponseProxyInfo struct {
	IP     string `json:"ip"`
	Status string `json:"status"`
	Info   string `json:"info"`
}

func proxyCheckerResponse(ip string, info string, status string) ResponseProxyInfo {
	return ResponseProxyInfo{Status: status, Info: info, IP: ip}
}

func (p *Handler) Check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.check"

		ip := strings.Split(r.RemoteAddr, ":")[0]

		err := isProxyHeader(r.Header)
		if err != nil {
			helpers.JSON(w, proxyCheckerResponse(ip, err.Error(), statusProxy), http.StatusBadRequest)
			return
		}

		helpers.JSON(w, proxyCheckerResponse(ip, "", statusNonProxy), http.StatusOK)
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
