package cart

import (
	"errors"
	"fmt"
	"interview/entities"

	"strconv"

	"gorm.io/gorm"
)

var itemPriceMapping = map[string]float64{
	"shoe":  100,
	"purse": 200,
	"bag":   300,
	"watch": 300,
}

type CartItemForm struct {
	Product  string `form:"product"   binding:"required"`
	Quantity string `form:"quantity"  binding:"required"`
}

func AddItemToCart(db *gorm.DB, product, quantity, cookie string) (string, bool) {
	var isCartNew bool
	var cartEntity entities.CartEntity
	result := db.Where(fmt.Sprintf("status = '%s' AND session_id = '%s'", entities.CartOpen, cookie)).First(&cartEntity)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", false
		}
		isCartNew = true
		cartEntity = entities.CartEntity{
			SessionID: cookie,
			Status:    entities.CartOpen,
		}
		db.Create(&cartEntity)
	}

	item, ok := itemPriceMapping[product]
	if !ok {
		return "invalid item name", false
	}

	intQuantity, err := strconv.ParseInt(quantity, 10, 0)
	if err != nil {
		return "invalid quantity", false
	}

	var cartItemEntity entities.CartItem
	if isCartNew {
		cartItemEntity = entities.CartItem{
			CartID:      cartEntity.ID,
			ProductName: product,
			Quantity:    int(intQuantity),
			Price:       item * float64(intQuantity),
		}
		db.Create(&cartItemEntity)
	} else {
		result = db.Where(" cart_id = ? and product_name  = ?", cartEntity.ID, product).First(&cartItemEntity)

		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return "", false
			}
			cartItemEntity = entities.CartItem{
				CartID:      cartEntity.ID,
				ProductName: product,
				Quantity:    int(intQuantity),
				Price:       item * float64(intQuantity),
			}
			db.Create(&cartItemEntity)

		} else {
			cartItemEntity.Quantity += int(intQuantity)
			cartItemEntity.Price += item * float64(intQuantity)
			db.Save(&cartItemEntity)
		}
	}

	return "", true
}
