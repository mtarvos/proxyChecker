package main

import (
	"log/slog"
	"proxyChecker/internal/app"
	"proxyChecker/internal/config"
	"proxyChecker/pkg/logging"
)

func main() {
	cfg := config.MustLoad()
	log := logging.InitLogger(cfg.Env)

	log.Info("Starting proxy checker", slog.String("env", cfg.Env))

	app.Run(log, cfg)
}
