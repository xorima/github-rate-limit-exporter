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

`docker pull ghcr.io/xorima/github-token-metrics:latest`
