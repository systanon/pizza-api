package main

import (
	"pizza-backend/internal/config"
	"pizza-backend/internal/database"
	"pizza-backend/internal/handler"
	"pizza-backend/internal/repository"
	"pizza-backend/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db := database.NewPostgres(cfg.DatabaseURL)
	defer db.Close()

	database.RunMigrations(cfg.DatabaseURL)

	// repositories
	productRepo := repository.NewProductRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	addonRepo := repository.NewAddonRepo(db)
	cartRepo := repository.NewCartRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// services
	productService := service.NewProductService(productRepo)
	addonService := service.NewAddonService(addonRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo)

	// handlers
	productHandler := handler.NewProductHandler(productService)
	addonHandler := handler.NewAddonHandler(addonService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)

	// router
	r := NewRouter(productHandler, categoryHandler, addonHandler, cartHandler, orderHandler)
	r.Run(cfg.Addr)
}
