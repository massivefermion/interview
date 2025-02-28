package router

import (
	"errors"
	"interview/cart"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateRemoveItemHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("ice_session_id")

		if err != nil || errors.Is(err, http.ErrNoCookie) || (cookie != nil && cookie.Value == "") {
			c.Redirect(302, "/")
			return
		}

		cart.RemoveItem(db, c.Query("cart_item_id"), cookie.Value)
		c.Redirect(302, "/")
	}
}
