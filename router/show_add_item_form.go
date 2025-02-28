package router

import (
	"errors"
	"fmt"
	"interview/cart"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateShowAddItemFormHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, err := c.Request.Cookie("ice_session_id")
		if errors.Is(err, http.ErrNoCookie) {
			c.SetCookie("ice_session_id", time.Now().String(), 3600, "/", "localhost", false, true)
		}

		data := map[string]interface{}{
			"Error": c.Query("error"),
			//"cartItems": cartItems,
		}

		cookie, err := c.Request.Cookie("ice_session_id")
		if err == nil {
			data["CartItems"] = cart.GetCartItemData(db, cookie.Value)
		}

		html, err := renderTemplate(data)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(200, html)
	}
}

func renderTemplate(pageData interface{}) (string, error) {
	// Read and parse the HTML template file
	tmpl, err := template.ParseFiles("./static/add_item_form.html")
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v ", err)
	}

	// Create a strings.Builder to store the rendered template
	var renderedTemplate strings.Builder

	err = tmpl.Execute(&renderedTemplate, pageData)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v ", err)
	}

	// Convert the rendered template to a string
	resultString := renderedTemplate.String()

	return resultString, nil
}
