package main

import (
	"context"
	"github.com/xorima/github-token-metrics/internal/app"
	"github.com/xorima/github-token-metrics/pkg/logger"
	"os"
)

func main() {
	l := logger.NewLogger(logger.ModeJSON, nil, os.Stdout)
	err := app.NewApp(l, os.Getenv("GITHUB_TOKEN")).Run(context.Background())
	if err != nil {
		l.Error(err.Error())
	}
}
