package cart

import (
	"fmt"
	"interview/entities"

	"gorm.io/gorm"
)

func GetCartItemData(db *gorm.DB, sessionID string) (items []map[string]interface{}) {
	var cartEntity entities.CartEntity
	result := db.Where(fmt.Sprintf("status = '%s' AND session_id = '%s'", entities.CartOpen, sessionID)).First(&cartEntity)
	if result.Error != nil {
		return
	}

	var cartItems []entities.CartItem
	result = db.Where(fmt.Sprintf("cart_id = %d", cartEntity.ID)).Find(&cartItems)
	if result.Error != nil {
		return
	}

	for _, cartItem := range cartItems {
		item := map[string]interface{}{
			"ID":       cartItem.ID,
			"Quantity": cartItem.Quantity,
			"Price":    cartItem.Price,
			"Product":  cartItem.ProductName,
		}

		items = append(items, item)
	}

	return items
}
