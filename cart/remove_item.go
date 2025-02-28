package cart

import (
	"fmt"
	"interview/entities"
	"strconv"

	"gorm.io/gorm"
)

func RemoveItem(db *gorm.DB, cartItemIDString string, cookie string) {
	if cartItemIDString == "" {
		return
	}

	var cartEntity entities.CartEntity
	result := db.Where(fmt.Sprintf("status = '%s' AND session_id = '%s'", entities.CartOpen, cookie)).First(&cartEntity)
	if result.Error != nil {
		return
	}

	if cartEntity.Status == entities.CartClosed {
		return
	}

	cartItemID, err := strconv.Atoi(cartItemIDString)
	if err != nil {
		return
	}

	var cartItemEntity entities.CartItem
	result = db.Where("id = ?", cartItemID).First(&cartItemEntity)
	if result.Error != nil {
		return
	}

	db.Delete(&cartItemEntity)
}
