package config

import "os"

type Config struct {
	DatabaseURL         string
	Addr                string
	StripeSecretKey     string
	StripeWebhookSecret string
	PaymentSuccessURL   string
	PaymentCancelURL    string
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
