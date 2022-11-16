package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "lagom_requests_total",
	Help: "The total number of requests",
})

func Report() {
	totalRequests.Inc()
}
