package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Serve func() error
}

func Create(db *gorm.DB, generateHandler func(*gorm.DB) *gin.Engine) App {
	return App{
		Serve: func() error {
			server := &http.Server{
				Addr:    ":8088",
				Handler: generateHandler(db),
			}

			return server.ListenAndServe()
		},
	}
}
