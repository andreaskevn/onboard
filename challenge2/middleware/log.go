package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		log.Println("Api Start",
			r.Method,
			r.RequestURI,
		)

		next.ServeHTTP(rw, r)

		log.Println("Api End",
			r.Method,
			r.RequestURI,
			"status:", rw.statusCode,
			"duration:", time.Since(start),
		)
	})
}