package app

import (
	"context"
	"github.com/google/go-github/v62/github"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	log         *slog.Logger
	client      *github.Client
	githubToken string
}

func NewApp(log *slog.Logger, githubToken string) *App {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &App{
		log:    log,
		client: github.NewClient(tc),
	}
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("Starting github token metrics")
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		a.log.Info("Starting server", slog.Int("port", 2112))
		if err := http.ListenAndServe(":2112", nil); err != nil {
			a.log.Error("Failed to start server", "error", err)
		}
	}()
	dur := 1 * time.Minute
	a.log.Info("Starting rate limit check", slog.Float64("interval_seconds", dur.Seconds()))
	a.process(ctx)
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.process(ctx)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *App) process(ctx context.Context) {
	if err := a.checkRateLimit(ctx); err != nil {
		a.log.Error("Failed to check rate limit", "error", err)
	}
}

func (a *App) getRateLimit(ctx context.Context) (*github.RateLimits, *github.Response, error) {
	a.log.Info("Checking Rate Limit")
	defer measureDuration("rate-limit", "get")()
	return a.client.RateLimit.Get(ctx)
}

func (a *App) checkRateLimit(ctx context.Context) error {
	rl, resp, err := a.getRateLimit(ctx)
	if err != nil {
		a.log.Error("Failed to get rate limit", "error", err)
		return err
	}
	setLastRunTime()
	setRateLimitMetrics(rl.GetCore(), "core")
	setRateLimitMetrics(rl.GetSearch(), "search")
	setRateLimitMetrics(rl.GetGraphQL(), "graphql")
	setRateLimitMetrics(rl.GetIntegrationManifest(), "integration_manifest")
	setRateLimitMetrics(rl.GetSourceImport(), "source_import")
	setRateLimitMetrics(rl.GetCodeScanningUpload(), "code_scanning_upload")
	setRateLimitMetrics(rl.GetActionsRunnerRegistration(), "actions_runner_registration")
	setRateLimitMetrics(rl.GetSCIM(), "scim")
	setRateLimitMetrics(rl.GetDependencySnapshots(), "dependency_snapshots")
	setRateLimitMetrics(rl.GetCodeSearch(), "code_search")
	setRateLimitMetrics(rl.GetAuditLog(), "audit_log")
	v := resp.Header.Get("github-authentication-token-expiration")
	if v != "" {
		expiry, err := time.Parse("2006-01-02 15:04:05 MST", v)
		if err != nil {
			a.log.Error("Failed to parse token expiry", "error", err)
			return err
		}
		setPatTokenExpiry(expiry.UnixMilli())

	}
	return nil
}
