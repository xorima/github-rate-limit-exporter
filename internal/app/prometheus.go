package app

import (
	"github.com/google/go-github/v71/github"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const namespace, subSystem = "github", "rate_limit"

var (
	rateLimitGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "limit"),
			Help: "The limit for different types of GitHub API requests",
		},
		[]string{"resource"},
	)
	rateRemainingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "remaining"),
			Help: "The remaining rate for different types of GitHub API requests",
		},
		[]string{"resource"},
	)
	rateResetGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "reset"),
			Help: "The reset time for different types of GitHub API requests",
		},
		[]string{"resource"},
	)
	patTokenExpiryGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "pat_token_expiry"),
			Help: "The expiry time for current token in ms",
		},
	)
	lastRunTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "last_run_time"),
			Help: "The last time the batch process checked for metrics",
		},
	)
	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation", "action"})
)

func init() {
	prometheus.MustRegister(rateLimitGauge, rateRemainingGauge, rateResetGauge, patTokenExpiryGauge, lastRunTime, httpDuration)
}

func setRateLimitMetrics(rate *github.Rate, name string) {
	rateLimitGauge.WithLabelValues(name).Set(float64(rate.Limit))
	rateRemainingGauge.WithLabelValues(name).Set(float64(rate.Remaining))
	rateResetGauge.WithLabelValues(name).Set(float64(rate.Reset.Unix()))
}
func setLastRunTime() {
	lastRunTime.SetToCurrentTime()
}

func setPatTokenExpiry(epoch int64) {
	patTokenExpiryGauge.Set(float64(epoch))
}

func measureDuration(operation, action string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues(operation, action).Observe(duration)
	}
}
