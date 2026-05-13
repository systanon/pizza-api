package config

import "os"

type Config struct {
	DatabaseURL string
	Addr        string
}

func Load() Config {
	addr := os.Getenv("PIZZA_ADDR")
	if addr == "" {
		addr = ":3002"
	}

	return Config{
		DatabaseURL: os.Getenv("PIZZA_DATABASE_URL"),
		Addr:        addr,
	}
}
