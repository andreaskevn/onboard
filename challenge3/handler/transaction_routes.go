package handler

import (
	"challenge3/server"
	"net/http"
)

func (t *TransactionHandler) MapRoutes() {
	t.mux.HandleFunc(
		server.NewAPIPath(http.MethodPost, "/transfer"),
		t.Transfer(),
	)
}
