package handler

import (
	"challenge3/server"
	"net/http"
	// "github.com/go-chi/chi/v5"
)

func (t *BankHandler) MapRoutes() {
	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodGet, "/banks"),
		t.GetBank(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodGet, "/banks/"),
		t.GetBankById(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodPost, "/banks"),
		t.CreateAcc(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodPut, "/banks/"),
		t.UpdateBank(),
	)

	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodDelete, "/banks/"),
		t.DeleteBank(),
	)
}
