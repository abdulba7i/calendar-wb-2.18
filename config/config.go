package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
	User        string        `yaml:"user" env:"HTTP_USER"`
	Password    string        `yaml:"password" env:"HTTP_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml" // дефолт
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	// Пароль и логин скрываем от пользователя
	log.Printf("Loaded config from %s: {HTTPServer: {Address: %s, Timeout: %s, IdleTimeout: %s, User: %s}}",
		configPath, cfg.HTTPServer.Address, cfg.HTTPServer.Timeout,
		cfg.HTTPServer.IdleTimeout, cfg.HTTPServer.User)

	return &cfg
}
