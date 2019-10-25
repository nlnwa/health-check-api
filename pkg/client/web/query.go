package web

import (
	"context"
	"net/http"
)

type Query interface {
	CheckVeidemannDashboard(ctx context.Context) (int, string, error)
}

func (ac Client) CheckVeidemannDashboard(ctx context.Context) (int, string, error) {
	resp, err := http.Head(ac.veidemannDashboardUrl)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, resp.Status, err
}
