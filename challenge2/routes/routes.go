package routes

import (
	"challenge2/handler"
	"net/http"
)

func RegisterItemRoutes(itemHandler *handler.ItemHandler) {

	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			itemHandler.GetItems(w, r)

		case http.MethodPost:
			itemHandler.CreateItem(w, r)

		case http.MethodDelete:
			itemHandler.DeleteItem(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

	})
}