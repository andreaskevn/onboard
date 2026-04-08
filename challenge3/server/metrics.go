package server

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var Meter = otel.Meter("bank-service")

var (
	HttpRequestCounter metric.Int64Counter
	HttpRequestDuration metric.Float64Histogram
)

var (
	AccountCreatedCounter metric.Int64Counter
	AccountFailedCounter  metric.Int64Counter
)

func InitMetrics() {
	var err error

	// HTTP
	HttpRequestCounter, err = Meter.Int64Counter("http_requests_total")
	if err != nil {
		panic(err)
	}

	HttpRequestDuration, err = Meter.Float64Histogram("http_request_duration_seconds")
	if err != nil {
		panic(err)
	}

	// Business
	AccountCreatedCounter, err = Meter.Int64Counter("account_created_total")
	if err != nil {
		panic(err)
	}

	AccountFailedCounter, err = Meter.Int64Counter("account_failed_total")
	if err != nil {
		panic(err)
	}
}