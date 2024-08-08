package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IP           string `env:"IP" env-default:"0.0.0.0"`
	HttpPort     uint   `env:"HTTP_PORT" env-default:"80"`
	HttpsPort    uint   `env:"HTTPS_PORT" env-default:"443"`
	KeyFile      string `env:"KEY" env-default:"server.crt"`
	CertFile     string `env:"CERT" env-default:"server.key"`
	HttpConnect  string
	HttpsConnect string
}

func Must() *Config {
	cfg := Config{}
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("env variable not found: %v", err)
		}
	}
	cfg.HttpConnect = fmt.Sprintf("%v:%v", cfg.IP, cfg.HttpPort)
	cfg.HttpsConnect = fmt.Sprintf("%v:%v", cfg.IP, cfg.HttpsPort)
	return &cfg
}
