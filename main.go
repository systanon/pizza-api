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

	// services
	productService := service.NewProductService(productRepo)
	addonService := service.NewAddonService(addonRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// handlers
	productHandler := handler.NewProductHandler(productService)
	addonHandler := handler.NewAddonHandler(addonService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// router
	r := NewRouter(productHandler, categoryHandler, addonHandler)
	r.Run(cfg.Addr)
}
