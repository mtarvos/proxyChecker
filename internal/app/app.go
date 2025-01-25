package app

import (
	"log/slog"
	"net/http"
	"proxyChecker/internal/adapters/external"
	"proxyChecker/internal/adapters/repository/sqlite"
	"proxyChecker/internal/config"
	"proxyChecker/internal/controller/http"
	"proxyChecker/internal/controller/http/handler"
	"proxyChecker/internal/controller/http/middleware"
	"proxyChecker/internal/service"
)

func Run(log *slog.Logger, cfg *config.Config) {
	storage, err := sqlite.New(cfg.StoragePath, log)
	if err != nil {
		log.Error("Failed to init storage", "error", err.Error())
	}

	proxyProvider := external.NewProxyApiClient(log)
	updater := service.NewUpdaterService(log, cfg.ProxyUpdateURL, proxyProvider, storage)
	updater.StartUpdateProxyThread()

	statsService := service.NewStatsService(log, storage)
	proxyService := service.NewProxy(storage)

	handler := handler.NewHandler(log, proxyService, statsService)

	r := router.New(handler)

	m := middleware.NewMiddleware(log)
	stack := m.GetStack()

	server := http.Server{
		Addr:    cfg.Address,
		Handler: stack(r),
	}

	log.Info("Server listening", slog.String("addr", cfg.Address))
	server.ListenAndServe()
}
