package calculator

import (
	"bytes"
	"fmt"
	"interview/pkg/db"
	"interview/pkg/entity"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var items = []string{"shoe", "purse", "bag", "watch"}

var engine *gin.Engine
var cookie *http.Cookie
var database *gorm.DB
var itemID uint

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	database = db.GetDatabase()

	engine = gin.Default()
	engine.POST("/add-item", AddItemToCart)
	engine.GET("/remove-cart-item", DeleteCartItem)

	cookie = &http.Cookie{
		Name:     "ice_session_id",
		Value:    time.Now().String(),
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
	}

	os.Exit(m.Run())
}

func TestAddItemToCart(t *testing.T) {
	product := items[rand.Intn(len(items))]
	quantity := rand.Intn(16) + 1

	body := url.Values{}
	body.Set("product", product)
	body.Set("quantity", fmt.Sprintf("%d", quantity))

	request, err := http.NewRequest("POST", "/add-item", bytes.NewBufferString(body.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	var cart entity.CartEntity
	result := database.Where("status = ? and session_id = ?", entity.CartOpen, cookie.Value).First(&cart)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	var cartItem entity.CartItem
	result = database.Where("cart_id = ?", cart.ID).First(&cartItem)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	itemID = cartItem.ID
	assert.Equal(t, cartItem.ProductName, product)
	assert.Equal(t, cartItem.Quantity, quantity)
}

func TestDeleteCartItem(t *testing.T) {
	request, err := http.NewRequest("GET", fmt.Sprintf("/remove-cart-item?cart_item_id=%d", itemID), nil)
	request.AddCookie(cookie)

	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	var cartItem entity.CartItem
	result := database.Where("id = ?", itemID).First(&cartItem)

	assert.NotNil(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}
