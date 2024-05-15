# github-token-metrics

A simple app to check the rate limits (and expiry) of the given github token and publish as prom metrics 

## Running

```bash
go run main.go
```

## Endpoints

`<host>:2112/metrics`

## Env Vars

`GITHUB_TOKEN` - the github token you want to check.

## Docker Image

`docker pull ghcr.io/xorima/github-token-metrics/github-token-metrics:latest`

## Example Metrics

```prometheus
# HELP github_pat_token_expiry The expirty time for current token
# TYPE github_pat_token_expiry gauge
github_pat_token_expiry 1.72001706e+09
# HELP github_rate_limit The limit for different types of GitHub API requests
# TYPE github_rate_limit gauge
github_rate_limit{resource="actions_runner_registration"} 10000
github_rate_limit{resource="audit_log"} 1750
github_rate_limit{resource="code_scanning_upload"} 1000
github_rate_limit{resource="code_search"} 10
github_rate_limit{resource="core"} 5000
github_rate_limit{resource="dependency_snapshots"} 100
github_rate_limit{resource="graphql"} 5000
github_rate_limit{resource="integration_manifest"} 5000
github_rate_limit{resource="scim"} 15000
github_rate_limit{resource="search"} 30
github_rate_limit{resource="source_import"} 100
# HELP github_rate_remaining The remaining rate for different types of GitHub API requests
# TYPE github_rate_remaining gauge
github_rate_remaining{resource="actions_runner_registration"} 10000
github_rate_remaining{resource="audit_log"} 1750
github_rate_remaining{resource="code_scanning_upload"} 1000
github_rate_remaining{resource="code_search"} 10
github_rate_remaining{resource="core"} 5000
github_rate_remaining{resource="dependency_snapshots"} 100
github_rate_remaining{resource="graphql"} 5000
github_rate_remaining{resource="integration_manifest"} 5000
github_rate_remaining{resource="scim"} 15000
github_rate_remaining{resource="search"} 30
github_rate_remaining{resource="source_import"} 100
# HELP github_rate_reset The reset time for different types of GitHub API requests
# TYPE github_rate_reset gauge
github_rate_reset{resource="actions_runner_registration"} 1.715769529e+09
github_rate_reset{resource="audit_log"} 1.715769529e+09
github_rate_reset{resource="code_scanning_upload"} 1.715769529e+09
github_rate_reset{resource="code_search"} 1.715765989e+09
github_rate_reset{resource="core"} 1.715769529e+09
github_rate_reset{resource="dependency_snapshots"} 1.715765989e+09
github_rate_reset{resource="graphql"} 1.715769529e+09
github_rate_reset{resource="integration_manifest"} 1.715769529e+09
github_rate_reset{resource="scim"} 1.715769529e+09
github_rate_reset{resource="search"} 1.715765989e+09
github_rate_reset{resource="source_import"} 1.715765989e+09
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 4.1297e-05
go_gc_duration_seconds{quantile="0.25"} 9.3405e-05
go_gc_duration_seconds{quantile="0.5"} 0.000123964
go_gc_duration_seconds{quantile="0.75"} 0.000187395
go_gc_duration_seconds{quantile="1"} 0.001793717
go_gc_duration_seconds_sum 0.012006006
go_gc_duration_seconds_count 56
```
