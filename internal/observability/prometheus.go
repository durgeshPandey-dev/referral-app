package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusHTTPHandler() http.Handler {
	return promhttp.HandlerFor(PrometheusRegistry(), promhttp.HandlerOpts{})
}
