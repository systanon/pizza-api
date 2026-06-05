package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL         string
	Addr                string
	StripeSecretKey     string
	StripeWebhookSecret string
	PaymentSuccessURL   string
	PaymentCancelURL    string
}

func (c Config) Validate() {
	required := map[string]string{
		"PIZZA_DATABASE_URL":  c.DatabaseURL,
		"STRIPE_SECRET_KEY":   c.StripeSecretKey,
		"STRIPE_WEBHOOK_SECRET": c.StripeWebhookSecret,
		"PAYMENT_SUCCESS_URL": c.PaymentSuccessURL,
		"PAYMENT_CANCEL_URL":  c.PaymentCancelURL,
	}
	for name, val := range required {
		if val == "" {
			log.Fatalf("required environment variable %s is not set", name)
		}
	}
}

func Load() Config {
	addr := os.Getenv("PIZZA_ADDR")
	if addr == "" {
		addr = ":3002"
	}

	return Config{
		DatabaseURL:         os.Getenv("PIZZA_DATABASE_URL"),
		Addr:                addr,
		StripeSecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		PaymentSuccessURL:   os.Getenv("PAYMENT_SUCCESS_URL"),
		PaymentCancelURL:    os.Getenv("PAYMENT_CANCEL_URL"),
	}
}
