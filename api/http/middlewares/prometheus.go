package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
	)

	responseDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(responseDuration)
}

func TraceNumberRequestAndTimeResponse(c *gin.Context) {
	start := time.Now()

	// Proceed to the next handler
	c.Next()

	// After handler processing
	requestCount.Inc()
	duration := time.Since(start).Seconds()
	responseDuration.Observe(duration)
}
