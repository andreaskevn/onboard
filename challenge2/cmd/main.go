package main

import (
	"log"
	"net/http"

	"challenge2/config"
	"challenge2/handler"
	"challenge2/repository"
	"challenge2/routes"
	"challenge2/service"
)

func main() {
	db, err := config.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database berhasil connect")

	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	itemHandler := handler.NewItemHandler(itemService)

	routes.RegisterItemRoutes(itemHandler)

	log.Println("server running on port 8080")

	http.ListenAndServe(":8000", nil)
}
