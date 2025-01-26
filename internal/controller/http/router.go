package router

import (
	"net/http"
)

type Handler interface {
	Proxy() http.HandlerFunc
	Stats() http.HandlerFunc
	Check() http.HandlerFunc
}

func New(handler Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /get-proxy", handler.Proxy())
	router.HandleFunc("GET /stats", handler.Stats())
	router.HandleFunc("GET /check", handler.Check())

	return router
}
