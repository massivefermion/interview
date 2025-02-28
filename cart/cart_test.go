package cart

import (
	"fmt"
	"interview/database"
	"interview/entities"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var items = []string{"shoe", "purse", "bag", "watch"}

var cookie string
var db *gorm.DB
var itemName string
var itemID uint

func TestMain(m *testing.M) {
	db = database.GetTestDatabase()
	cookie = time.Now().String()
	os.Exit(m.Run())
}

func TestAddItemToCart(t *testing.T) {
	product := items[rand.Intn(len(items))]
	quantity := rand.Intn(16) + 1

	body := url.Values{}
	body.Set("product", product)
	body.Set("quantity", fmt.Sprintf("%d", quantity))

	AddItemToCart(db, product, fmt.Sprintf("%d", quantity), cookie)
	var cart entities.CartEntity
	result := db.Where("status = ? and session_id = ?", entities.CartOpen, cookie).First(&cart)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	var cartItem entities.CartItem
	result = db.Where("cart_id = ?", cart.ID).First(&cartItem)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	itemName = cartItem.ProductName
	itemID = cartItem.ID

	assert.Equal(t, cartItem.ProductName, product)
	assert.Equal(t, cartItem.Quantity, quantity)
}

func TestGetCartItemData(t *testing.T) {
	items := GetCartItemData(db, cookie)
	item := items[0]

	assert.Equal(t, itemName, item["Product"])
}

func TestDeleteCartItem(t *testing.T) {
	RemoveItem(db, fmt.Sprintf("%d", itemID), cookie)
	var cartItem entities.CartItem
	result := db.Where("id = ?", itemID).First(&cartItem)

	assert.NotNil(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}
