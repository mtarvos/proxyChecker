package main

import (
	"log/slog"
	"proxyChecker/internal/app"
	"proxyChecker/internal/config"
	"proxyChecker/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.InitLogger(cfg.Env, cfg.LogFile)

	log.Info("Starting proxy checker", slog.String("env", cfg.Env))

	app.Run(log, cfg)
}
