package app

import (
	"log/slog"
	"net/http"
	"os"
	"proxyChecker/internal/adapter/external"
	"proxyChecker/internal/adapter/repository/sqlite"
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
		os.Exit(1)
		return
	}

	err = storage.MigrationsUP()
	if err != nil {
		log.Error("Failed to set migrations", "error", err.Error())
		os.Exit(1)
	}

	checkerClient := external.NewAbstractApiClient(log, cfg.ProxyCheckerURL, cfg.ProxyType)
	checkService := service.NewCheckerService(log, storage, checkerClient)
	go checkService.StartCheckerRoutine(cfg.CheckRoutineCount)

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
