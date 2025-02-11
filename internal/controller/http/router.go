package router

import (
	"net/http"
)

type Handler interface {
	Proxy() http.HandlerFunc
	Stats() http.HandlerFunc
	Check() http.HandlerFunc
	Next() http.HandlerFunc
}

func New(handler Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /proxy", handler.Proxy())
	router.HandleFunc("GET /stats", handler.Stats())
	router.HandleFunc("GET /check", handler.Check())
	router.HandleFunc("GET /next", handler.Next())

	return router
}
