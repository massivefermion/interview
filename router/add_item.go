package router

import (
	"errors"
	"fmt"
	"interview/cart"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type CartItemForm struct {
	Product  string `form:"product"   binding:"required"`
	Quantity string `form:"quantity"  binding:"required"`
}

func generateAddItemHandler(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("ice_session_id")

		if err != nil || errors.Is(err, http.ErrNoCookie) || (cookie != nil && cookie.Value == "") {
			c.Redirect(302, "/")
			return
		}

		cartItemForm, err := getCartItemForm(c)
		if err != nil {
			c.Redirect(302, "/?error="+err.Error())
			return
		}

		msg, _ := cart.AddItemToCart(db, cartItemForm.Product, cartItemForm.Quantity, cookie.Value)
		c.Redirect(302, "/"+msg)
	}
}

func getCartItemForm(c *gin.Context) (*CartItemForm, error) {
	if c.Request.Body == nil {
		return nil, fmt.Errorf("body cannot be nil")
	}

	form := &CartItemForm{}

	if err := binding.FormPost.Bind(c.Request, form); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return form, nil
}
