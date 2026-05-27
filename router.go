package main

import (
	"pizza-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	addonHandler *handler.AddonHandler,

) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/categories", categoryHandler.GetCategoryList)
		v1.GET("/categories/:id/addons", addonHandler.GetAddonListByCategoryID)

		v1.GET("/products", productHandler.GetProductList)
		v1.GET("/products/:id", productHandler.GetProductByID)
	}
	r.NoRoute(handler.NotFoundHandler)

	return r
}
