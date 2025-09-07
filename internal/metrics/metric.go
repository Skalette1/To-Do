package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number of http requests",
		},
		[]string{"handler", "method", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "The duration of http requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler"},
	)
)

func Init() {
	prometheus.MustRegister(HttpRequestsTotal, HttpRequestDuration)
	http.Handle("/metrics", promhttp.Handler())
}

func IncRequests(handler, method, status string) {
	HttpRequestsTotal.WithLabelValues(handler, method, status).Inc()
}

func ObserveDuration(handler string, seconds float64) {
	HttpRequestDuration.WithLabelValues(handler).Observe(seconds)
}
