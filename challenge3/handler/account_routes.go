package handler

import (
	"challenge3/server"
	"net/http"
	// "github.com/go-chi/chi/v5"
)

func (t *AccountHandler) MapRoutes() {
	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodGet, "/accounts"),
		t.Get(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodGet, "/accounts/"),
		t.AccountRouter(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodPost, "/accounts"),
		t.CreateAcc(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodPut, "/accounts/"),
		t.UpdateAcc(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodDelete, "/accounts/"),
		t.DeleteAcc(),
	)

	// t.mux.HandleFunc(
	// 	server.NewAPIPath(http.MethodGet, "/accounts/:id/transfer/"),
	// 	t.GetHistory(),
	// )
}
