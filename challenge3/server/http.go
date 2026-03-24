package server

import (
	"encoding/json"
	"net/http"
	"challenge3/dto"
)

func ApplicationMiddlewareResponse(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func HandleRouteNotFound(mux *http.ServeMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h, pattern := mux.Handler(r)
		if pattern == "" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Not Found",
				Data: "Not found",
			})
			return
		}

		h.ServeHTTP(w, r)
	}
}