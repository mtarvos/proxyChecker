package app

import (
	"log/slog"
	"net/http"
	"os"
	"proxyChecker/internal/config"
	"proxyChecker/internal/controller/http"
	"proxyChecker/internal/controller/http/handler"
	"proxyChecker/internal/controller/http/middleware"
	"proxyChecker/internal/infrastructure/client"
	"proxyChecker/internal/infrastructure/repository/sqlite"
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

	proxyProvider := client.NewProxyProvider(log)
	updater := service.NewUpdaterService(log, cfg.ProxyUpdateURL, proxyProvider, storage)
	updater.StartUpdateProxyRoutine()

	statsService := service.NewStatsService(log, storage)

	checkerClient := client.NewChecker(log, cfg.CheckerURL, cfg.ProxyType)
	checkService := service.NewCheckerService(log, storage, checkerClient)
	go checkService.StartCheckerRoutine(cfg.CheckRoutineCount)

	infoClient := client.NewAbstractAPI(log, cfg.InfoURL, cfg.Key)
	infoService := service.NewInfoService(log, infoClient, storage)
	go infoService.StartInfoRoutine(cfg.InfoRoutineCount)

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
