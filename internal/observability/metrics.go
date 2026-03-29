package observability

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsOnce sync.Once
	registry    = prometheus.NewRegistry()
	enabled     bool

	httpRequestsTotal *prometheus.CounterVec
	httpDuration      *prometheus.HistogramVec
	httpInflight      prometheus.Gauge

	uploadRequestsTotal *prometheus.CounterVec
	queueDepth          prometheus.Gauge
	queueJobsTotal      *prometheus.CounterVec
	contactsTotal       *prometheus.CounterVec
	emailDuration       *prometheus.HistogramVec
)

func InitMetrics() {
	metricsOnce.Do(func() {
		httpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "onix_http_requests_total",
				Help: "Total inbound HTTP requests by method, route and status.",
			},
			[]string{"method", "route", "status"},
		)

		httpDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "onix_http_request_duration_seconds",
				Help:    "End-to-end HTTP request duration in seconds.",
				Buckets: []float64{0.01, 0.03, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 30},
			},
			[]string{"method", "route", "status"},
		)

		httpInflight = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "onix_http_inflight_requests",
				Help: "Current number of in-flight HTTP requests.",
			},
		)

		uploadRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "onix_upload_requests_total",
				Help: "Count of upload API requests by result.",
			},
			[]string{"result"},
		)

		queueDepth = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "onix_queue_depth",
				Help: "Current queue backlog depth.",
			},
		)

		queueJobsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "onix_queue_jobs_total",
				Help: "Queue job lifecycle events (enqueued, sent, failed, cancelled).",
			},
			[]string{"event"},
		)

		contactsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "onix_contacts_total",
				Help: "Total contacts processed by pipeline stage.",
			},
			[]string{"stage"},
		)

		emailDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "onix_email_send_duration_seconds",
				Help:    "Email API call duration in seconds by result.",
				Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10},
			},
			[]string{"result"},
		)

		registry.MustRegister(
			httpRequestsTotal,
			httpDuration,
			httpInflight,
			uploadRequestsTotal,
			queueDepth,
			queueJobsTotal,
			contactsTotal,
			emailDuration,
		)
		enabled = true
	})
}

func PrometheusRegistry() *prometheus.Registry {
	return registry
}

func IncInflightHTTPRequests() {
	if !enabled || httpInflight == nil {
		return
	}
	httpInflight.Inc()
}

func DecInflightHTTPRequests() {
	if !enabled || httpInflight == nil {
		return
	}
	httpInflight.Dec()
}

func ObserveHTTPRequest(method, route string, status int, duration time.Duration) {
	if !enabled || httpRequestsTotal == nil || httpDuration == nil {
		return
	}
	statusStr := strconv.Itoa(status)
	httpRequestsTotal.WithLabelValues(method, route, statusStr).Inc()
	httpDuration.WithLabelValues(method, route, statusStr).Observe(duration.Seconds())
}

func IncUploadRequest(result string) {
	if !enabled || uploadRequestsTotal == nil {
		return
	}
	uploadRequestsTotal.WithLabelValues(result).Inc()
}

func SetQueueDepth(depth int) {
	if !enabled || queueDepth == nil {
		return
	}
	queueDepth.Set(float64(depth))
}

func IncQueueJobEvent(event string) {
	if !enabled || queueJobsTotal == nil {
		return
	}
	queueJobsTotal.WithLabelValues(event).Inc()
}

func AddContacts(stage string, count int) {
	if !enabled || contactsTotal == nil {
		return
	}
	contactsTotal.WithLabelValues(stage).Add(float64(count))
}

func ObserveEmailDuration(duration time.Duration, result string) {
	if !enabled || emailDuration == nil {
		return
	}
	emailDuration.WithLabelValues(result).Observe(duration.Seconds())
}
