package web

import (
	"context"
	"fmt"
	"net/http"
)

type Query interface {
	CheckVeidemannDashboard(ctx context.Context) (int, string, error)
}

func (ac Client) CheckVeidemannDashboard(ctx context.Context) (int, string, error) {
	resp, err := http.Head(ac.veidemannDashboardUrl)
	if err != nil {
		return 0, "", fmt.Errorf("failed to request url: %w", err)
	}
	return resp.StatusCode, resp.Status, err
}
