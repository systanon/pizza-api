package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://pizza_user:pizza_pass@postgres:5432/pizza_db?sslmode=disable")
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("connected to database")

	select {}
}
