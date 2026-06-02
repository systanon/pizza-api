package main

import (
	"pizza-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	addonHandler *handler.AddonHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/categories", categoryHandler.GetCategoryList)
		v1.GET("/categories/:id/addons", addonHandler.GetAddonListByCategoryID)

		v1.GET("/products", productHandler.GetProductList)
		v1.GET("/products/:id", productHandler.GetProductByID)

		v1.POST("/cart", cartHandler.CreateCart)
		v1.GET("/cart/:id", cartHandler.GetCart)
		v1.DELETE("/cart/:id", cartHandler.ClearCart)
		v1.POST("/cart/:id/items", cartHandler.AddItem)
		v1.PATCH("/cart/:id/items/:itemId", cartHandler.UpdateItem)
		v1.DELETE("/cart/:id/items/:itemId", cartHandler.RemoveItem)

		v1.POST("/orders", orderHandler.CreateOrder)
		v1.GET("/orders/:id", orderHandler.GetOrder)
	}
	r.NoRoute(handler.NotFoundHandler)

	return r
}
