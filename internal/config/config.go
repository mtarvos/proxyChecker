package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"proxyChecker/internal/entity"
	"time"
)

type Config struct {
	Env            string `yaml:"env" env-default:"local"`
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	HTTPServer     `yaml:"http_server" env-required:"true"`
	ProxyUpdateURL string `yaml:"proxy_update_url" env-required:"true"`
	Checker        `yaml:"proxy_checker" env-required:"true"`
	AbstractAPI    `yaml:"abstract_api" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Checker struct {
	CheckerURL        string        `yaml:"checker_url" env-required:"true"`
	ProxyType         entity.Status `yaml:"proxy_type" env-required:"true"`
	CheckRoutineCount int           `yaml:"routine_count" env-required:"true"`
}

type AbstractAPI struct {
	InfoURL          string `yaml:"info_url" env-required:"true"`
	Key              string `yaml:"key" env-required:"true"`
	InfoRoutineCount int    `yaml:"routine_count" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Can not read config: %s", err)
	}

	return &cfg
}
