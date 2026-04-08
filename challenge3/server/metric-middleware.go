package server

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		attrs := []attribute.KeyValue{
			attribute.String("method", r.Method),
			attribute.String("route", r.URL.Path),
		}

		HttpRequestCounter.Add(
			r.Context(),
			1,
			metric.WithAttributes(attrs...),
		)

		HttpRequestDuration.Record(
			r.Context(),
			duration,
			metric.WithAttributes(attrs...),
		)
	})
}