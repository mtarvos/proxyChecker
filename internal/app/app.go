package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"proxyChecker/internal/config"
	"proxyChecker/internal/controller/http"
	"proxyChecker/internal/controller/http/handler"
	"proxyChecker/internal/controller/http/middleware"
	"proxyChecker/internal/infrastructure/client"
	"proxyChecker/internal/infrastructure/repository/sqlite"
	"proxyChecker/internal/service"
	"proxyChecker/pkg/logging"
	"sync"
	"time"
)

func Run(log *slog.Logger, cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logging.ContextWithLogger(ctx, log)

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

	var wg sync.WaitGroup

	proxyProvider := client.NewProxyProvider()
	updater := service.NewUpdaterService(cfg.ProxyUpdateURL, proxyProvider, storage)
	wg.Add(1)
	updater.StartUpdateProxyRoutine(ctx, &wg)

	statsService := service.NewStatsService(storage)

	checkerClient := client.NewChecker(cfg.CheckerURL, cfg.ProxyType)
	checkService := service.NewCheckerService(storage, checkerClient)
	wg.Add(1)
	go checkService.StartCheckerRoutine(ctx, cfg.CheckRoutineCount, &wg)

	infoClient := client.NewAbstractAPI(cfg.InfoURL, cfg.Key)
	infoService := service.NewInfoService(infoClient, storage)
	wg.Add(1)
	go infoService.StartInfoRoutine(ctx, cfg.InfoRoutineCount, &wg)

	nextService := service.NewNextService(storage)
	proxyService := service.NewProxy(storage)
	handler := handler.NewHandler(proxyService, statsService, nextService)

	r := router.New(handler)

	m := middleware.NewMiddleware(log)
	stack := m.GetStack()

	server := http.Server{
		Addr:    cfg.Address,
		Handler: stack(r),
	}

	go func() {
		log.Info("Server listening", slog.String("addr", cfg.Address))
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	sigReceived := <-signalChan
	log.Info("Received signal:", slog.String("os.Signal", sigReceived.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Error("Error during server shutdown", slog.String("error", err.Error()))
	}
	log.Info("Server shutdown completed")

	cancel()
	wg.Wait()

	log.Info("All goroutines completed")
	log.Info("App gracefully stopped")
}
