package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	
)

// Main metrics we track
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"method", "path"},
	)

	successRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_success_total",
			Help: "Number of successful HTTP requests (2xx)",
		},
		[]string{"method", "path"},
	)

	failedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_failed_total",
			Help: "Number of failed HTTP requests (4xx/5xx)",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response durations per path",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	inFlightRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Current number of in-flight requests",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(
		totalRequests,
		successRequests,
		failedRequests,
		requestDuration,
		inFlightRequests,
	)
}

// PrometheusMiddleware tracks all core metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path // fallback for unknown routes
		}

		inFlightRequests.WithLabelValues(path).Inc()
		defer inFlightRequests.WithLabelValues(path).Dec()

		// Continue processing
		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := c.Writer.Status()
		status := fmt.Sprintf("%d", statusCode)

		// Record metrics
		totalRequests.WithLabelValues(c.Request.Method, path).Inc()
		requestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

		if statusCode >= 200 && statusCode < 300 {
			successRequests.WithLabelValues(c.Request.Method, path).Inc()
		} else {
			failedRequests.WithLabelValues(c.Request.Method, path, status).Inc()
		}
	}
}


