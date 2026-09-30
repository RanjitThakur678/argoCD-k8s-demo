package platform

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RequestDuration is labeled by domain (catalog/cart/reviews/checkout/payment) so Grafana can break
// down latency and cache-hit behavior per read/write-heavy path instead of one flat HTTP metric.
// Deliberately no raw path label - product/order IDs in the URL would blow up cardinality in Prometheus.
var RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "backend_request_duration_seconds",
	Help: "Request duration by domain, method and outcome.",
}, []string{"domain", "method", "status"})

var CacheResult = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "backend_cache_result_total",
	Help: "Cache hits/misses by domain - watch this to see if catalog/reviews caching is earning its keep.",
}, []string{"domain", "result"})

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
