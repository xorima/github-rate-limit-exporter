package app

import (
	"context"
	"github.com/google/go-github/v62/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"time"
)

const namespace, subSystem = "github", "rate_limit"

type App struct {
	log                 *slog.Logger
	client              *github.Client
	rateLimitGauge      *prometheus.GaugeVec
	rateRemainingGauge  *prometheus.GaugeVec
	rateResetGauge      *prometheus.GaugeVec
	patTokenExpiryGauge prometheus.Gauge
	lastRunTime         prometheus.Gauge
}

func NewApp(log *slog.Logger, githubToken string) *App {
	rateLimitGauge, rateRemainingGauge, rateResetGauge, patTokenExpiryGauge, lastRunTime := registerMetrics()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &App{
		log:                 log,
		client:              github.NewClient(tc),
		rateLimitGauge:      rateLimitGauge,
		rateRemainingGauge:  rateRemainingGauge,
		rateResetGauge:      rateResetGauge,
		patTokenExpiryGauge: patTokenExpiryGauge,
		lastRunTime:         lastRunTime,
	}
}

func registerMetrics() (*prometheus.GaugeVec, *prometheus.GaugeVec, *prometheus.GaugeVec, prometheus.Gauge, prometheus.Gauge) {

	rateLimitGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "limit"),
			Help: "The limit for different types of GitHub API requests",
		},
		[]string{"resource"},
	)
	rateRemainingGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "remaining"),
			Help: "The remaining rate for different types of GitHub API requests",
		},
		[]string{"resource"},
	)
	rateResetGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "reset"),
			Help: "The reset time for different types of GitHub API requests",
		},
		[]string{"resource"},
	)
	patTokenExpiryGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "pat_token_expiry"),
			Help: "The expiry time for current token",
		},
	)
	lastRunTime := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "last_run_time"),
			Help: "The last time the batch process checked for metrics",
		},
	)

	prometheus.MustRegister(rateLimitGauge, rateRemainingGauge, rateResetGauge, patTokenExpiryGauge, lastRunTime)
	return rateLimitGauge, rateRemainingGauge, rateResetGauge, patTokenExpiryGauge, lastRunTime
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("Running the app")
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		a.log.Info("Starting server", slog.Int("port", 2112))
		if err := http.ListenAndServe(":2112", nil); err != nil {
			a.log.Error("Failed to start server", "error", err)
		}
	}()
	dur := 1 * time.Minute
	a.log.Info("Starting rate limit check", slog.Float64("interval_seconds", dur.Seconds()))
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := a.checkRateLimit(ctx); err != nil {
				a.log.Error("Failed to check rate limit", "error", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *App) checkRateLimit(ctx context.Context) error {
	a.log.Info("Checking Rate Limit")
	rl, resp, err := a.client.RateLimit.Get(ctx)
	if err != nil {
		a.log.Error("Failed to get rate limit", "error", err)
		return err
	}
	a.lastRunTime.SetToCurrentTime()
	a.setRateLimitMetrics(rl.GetCore(), "core")
	a.setRateLimitMetrics(rl.GetSearch(), "search")
	a.setRateLimitMetrics(rl.GetGraphQL(), "graphql")
	a.setRateLimitMetrics(rl.GetIntegrationManifest(), "integration_manifest")
	a.setRateLimitMetrics(rl.GetSourceImport(), "source_import")
	a.setRateLimitMetrics(rl.GetCodeScanningUpload(), "code_scanning_upload")
	a.setRateLimitMetrics(rl.GetActionsRunnerRegistration(), "actions_runner_registration")
	a.setRateLimitMetrics(rl.GetSCIM(), "scim")
	a.setRateLimitMetrics(rl.GetDependencySnapshots(), "dependency_snapshots")
	a.setRateLimitMetrics(rl.GetCodeSearch(), "code_search")
	a.setRateLimitMetrics(rl.GetAuditLog(), "audit_log")
	v := resp.Header.Get("github-authentication-token-expiration")
	if v != "" {
		expiry, err := time.Parse("2006-01-02 15:04:05 MST", v)
		if err != nil {
			a.log.Error("Failed to parse token expiry", "error", err)
			return err
		}
		a.patTokenExpiryGauge.Set(float64(expiry.UnixMilli()))
	}
	return nil
}

func (a *App) setRateLimitMetrics(rate *github.Rate, name string) {
	a.rateLimitGauge.WithLabelValues(name).Set(float64(rate.Limit))
	a.rateRemainingGauge.WithLabelValues(name).Set(float64(rate.Remaining))
	a.rateResetGauge.WithLabelValues(name).Set(float64(rate.Reset.Unix()))

}
