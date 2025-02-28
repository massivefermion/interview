package main

import (
	"interview/app"
	"interview/database"
	"interview/router"
)

func main() {
	db := database.GetDatabase()
	app := app.Create(db, router.Create)
	app.Serve()
}
