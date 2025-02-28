package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(db *gorm.DB) *gin.Engine {
	handler := gin.Default()

	handler.GET("/", generateShowAddItemFormHandler(db))
	handler.POST("/add-item", generateAddItemHandler(db))
	handler.GET("/remove-cart-item", generateRemoveItemHandler(db))

	return handler
}
